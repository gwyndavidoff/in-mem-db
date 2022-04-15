package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/btree"
)

type DB struct {
	// Using btree for both mean O log n for retrieving, setting and deleting data
	// In the worst case, it also means the mem requirements are 2x if all Database values are distinct
	// Also, technically because of the implementation of btree, copying it for transactions would
	// not have doubled the space taken, but it was my understanding that using that would have gone
	// against the intention of that requirement, so I chose another route
	Database     *btree.BTree
	Counts       *btree.BTree
	Transactions []Transaction
}

type Transaction struct {
	Commands []string
}

//Node implements the btree Item Interface
type Node struct {
	key   string
	value string
}

func (n Node) Less(than btree.Item) bool {
	return n.key < than.(Node).key
}

//CountNode implements the btree Item Interface
type CountNode struct {
	key   string
	value int
}

func (cn CountNode) Less(than btree.Item) bool {
	return cn.key < than.(CountNode).key
}

func Init() *DB {
	return &DB{
		Database:     btree.New(2),
		Counts:       btree.New(2),
		Transactions: []Transaction{},
	}
}

// Handle uses a DB instance and an input string to perform all possible functions
// It's a bit naive at the moment, if an invalid command is issued it's ignored
// And if the number of args are greater than expected, they're ignored
// Only too few args are checked and return an error, Also, it is not case sensitive.
func (d *DB) Handle(input string) {
	args := strings.Split(input, " ")
	command := args[0]
	if strings.EqualFold(command, "END") {
		os.Exit(0)
	} else if strings.EqualFold(command, "SET") {
		if len(args) < 3 {
			fmt.Printf("SET requires 2 string inputs\n")
			return
		}
		d.Set(args[1], args[2])
	} else if strings.EqualFold(command, "GET") {
		if len(args) < 2 {
			fmt.Printf("GET requires 1 string inputs\n")
			return
		}
		fmt.Printf("%v\n", d.Get(args[1]))
	} else if strings.EqualFold(command, "DELETE") {
		if len(args) < 2 {
			fmt.Printf("DELETE requires 1 string inputs\n")
			return
		}
		d.Delete(args[1])
	} else if strings.EqualFold(command, "COUNT") {
		if len(args) < 2 {
			fmt.Printf("COUNT requires 1 string inputs\n")
			return
		}
		fmt.Println(d.Count(args[1]))
	} else if strings.EqualFold(command, "BEGIN") {
		d.Begin()
	} else if strings.EqualFold(command, "ROLLBACK") {
		d.Rollback()
	} else if strings.EqualFold(command, "COMMIT") {
		d.Commit()
	}
}

// A private function that processes the commands sent by a transaction rollback
func (d *DB) transactionRollbackHandle(input string) {
	args := strings.Split(input, " ")
	command := args[0]
	if strings.EqualFold(command, "SET") {
		d.transactionSet(args[1], args[2], true)
	} else if strings.EqualFold(command, "DELETE") {
		d.transactionDelete(args[1], true)
	}
}

// Set adds the new key/value pair to the Database
// increments the counter for that value, and
// decrements the counter if there was a previous value
func (d *DB) Set(key, value string) {
	d.transactionSet(key, value, false)
}

// A private version of Set so that during a transaction rollback
// transaction commands aren't added.
// rollback(bool) True: don't save undo command.
// False: save undo command (if there's an open transaction)
func (d *DB) transactionSet(key, value string, rollback bool) {
	n := Node{
		key:   key,
		value: value,
	}
	if oldNode := d.Database.ReplaceOrInsert(n); oldNode != nil {
		countOldNode := d.Counts.Get(CountNode{key: oldNode.(Node).value}).(CountNode)
		countOldNode.value--
		if countOldNode.value == 0 {
			d.Counts.Delete(countOldNode)
		} else {
			d.Counts.ReplaceOrInsert(countOldNode)
		}

		// If there's a transaction, store the "undo" - this is for a replace so it SETs the old key/value
		// Add it to the beginning of the list so it's performed in reverse order
		if len(d.Transactions) > 0 && !rollback {
			original := oldNode.(Node)
			lastTransaction := d.Transactions[len(d.Transactions)-1]
			lastTransaction.Commands = append([]string{fmt.Sprintf("SET %s %s", original.key, original.value)}, lastTransaction.Commands...)
			d.Transactions[len(d.Transactions)-1] = lastTransaction
		}
	} else {
		// If there's a transaction, store the "undo" - this is for an insert, so it DELETEs the new key
		// Add it to the beginning of the list so it's performed in reverse order
		if len(d.Transactions) > 0 && !rollback {
			lastTransaction := d.Transactions[len(d.Transactions)-1]
			lastTransaction.Commands = append([]string{fmt.Sprintf("DELETE %s", key)}, lastTransaction.Commands...)
			d.Transactions[len(d.Transactions)-1] = lastTransaction
		}
	}
	// Retrieve the current count of the value just saved
	countNode := CountNode{}
	if cn := d.Counts.Get(CountNode{key: n.value}); cn != nil {
		countNode = cn.(CountNode)
	}
	//Because int is initialized at 0 we can just set the key and ++ the value and then save, and not check for 0
	countNode.key = n.value
	countNode.value++
	d.Counts.ReplaceOrInsert(countNode)
}

// Get retrieves the matching value for the given key from the Database or "NULL"
func (d *DB) Get(key string) string {
	if n := d.Database.Get(Node{key: key}); n != nil {
		return n.(Node).value
	}
	return "NULL"
}

// Delete removes the key/value pair from the Database if it exists
// and decrements the count for the value, removing it from Counts ir it's now 0
func (d *DB) Delete(key string) {
	d.transactionDelete(key, false)
}

// A private version of Set so that during a transaction rollback
// transaction commands aren't added.
// rollback(bool) True: don't save undo command.
// False: save undo command (if there's an open transaction)
func (d *DB) transactionDelete(key string, rollback bool) {
	if d.Database.Len() == 0 {
		return
	}
	deleted := d.Database.Delete(Node{key: key})
	if deleted == nil {
		return
	}
	// If there's a transaction, store the "undo"
	if len(d.Transactions) > 0 && !rollback {
		lastTransaction := d.Transactions[len(d.Transactions)-1]
		deletedNode := deleted.(Node)
		lastTransaction.Commands = append([]string{fmt.Sprintf("SET %s %s", deletedNode.key, deletedNode.value)}, lastTransaction.Commands...)
		d.Transactions[len(d.Transactions)-1] = lastTransaction
	}
	// Being thorough with the nil check here...if the value was in the db it should be in Counts
	countNode := CountNode{}
	if cn := d.Counts.Get(CountNode{key: deleted.(Node).value}); cn != nil {
		countNode = cn.(CountNode)
		countNode.value--
		if countNode.value == 0 {
			d.Counts.Delete(countNode)
		} else {
			d.Counts.ReplaceOrInsert(countNode)
		}
	}
}

// Count returns the saved current count of the value from the Counts btree
func (d *DB) Count(value string) int {
	if n := d.Counts.Get(CountNode{key: value}); n != nil {
		return n.(CountNode).value
	}
	return 0
}

// Begin adds a new Transaction list item for storing undo Commands
func (d *DB) Begin() {
	d.Transactions = append(d.Transactions, Transaction{Commands: []string{}})
}

// Commit discards all open Transactions so that they are "permanent"
func (d *DB) Commit() {
	d.Transactions = nil
}

// Rollback takes the latest Transaction, performs the undo commands
// And then discards the rolled back Transaction
func (d *DB) Rollback() {
	if len(d.Transactions) == 0 {
		return
	}
	lastTransaction := d.Transactions[len(d.Transactions)-1]
	for _, cmd := range lastTransaction.Commands {
		d.transactionRollbackHandle(cmd)
	}
	d.Transactions = d.Transactions[0 : len(d.Transactions)-1]
}

// A private function primarily used to help while making changes and testing
func (d *DB) printTransactions() {
	fmt.Println(d.Transactions)
}
