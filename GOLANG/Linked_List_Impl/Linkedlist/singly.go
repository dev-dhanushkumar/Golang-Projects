package Linkedlist

import "fmt"

type node struct {
	data int
	next *node
}

type linkedList  struct {
	head *node 
}

func New() *linkedList {
	return &linkedList{}
}

func (l *linkedList) InsertAtBeginning(data int)  {
	newNode := &node {
		data: data,
		next: l.head,
	}
	l.head = newNode
}

func (l *linkedList)InsertAtMiddle(data int)  {
	newNode := &node {data: data}

	// case1: empty List
	if l.head == nil {
		l.head = newNode
		return
	}

	// Step 1: Get lenght of List
	len := l.GetLenghtofLinkedList()

	// Step 2: find middle value
	mid := len / 2 

	// Step 3: Traverse to middle node
	current := l.head
	for i := 1; i < mid; i++ {
		current = current.next
	}

	// Step 4: Insert the node
	newNode.next = current.next
	current.next = newNode
}

func (l *linkedList) InsertAtEnd(data int)  {
	newNode := &node {
		data: data,
	}

	if l.head == nil {
		l.head = newNode
		return
	}

	current := l.head
	for current.next != nil {
		current = current.next
	}

	current.next = newNode
}

// Inset Element at given position
func (l *linkedList)InsertAtIdx(data int, idx int)  {
	newNode := &node { data: data }

	// Empty list
	if l.head == nil {
		l.head = newNode
		return
	}

	// Get Length of Linked List
	len := l.GetLenghtofLinkedList()
	if idx > len {
		fmt.Println("Hay you given idx value not include list lenght please try again other idx value...")
		return
	}

	current := l.head
	for i := 1; i < idx; i++ {
		current = current.next
	}

	newNode.next = current.next
	current.next = newNode
}

func (l *linkedList) Display() {
	current := l.head
	for current != nil {
		fmt.Print(current.data, " -> ")
		current = current.next
	}
	fmt.Print("nil")
	fmt.Println()
}

/* 
 Implement the Get Lenght of Linked List function
 func (l *linkedList)GetLenghtofList() int {
 	**---- block of impl ---- **
 }
 */

 func (l *linkedList)GetLenghtofLinkedList() int {
	 len := 0
	 current := l.head

	 for current != nil {
		 len++
		 current = current.next
	 }

	 return len
 }

 /*
  * Reverse the linked list
 */
func (l *linkedList)ReverseList()  {
	var prev *node = nil
	current := l.head

	for current != nil {
		next := current.next
		current.next = prev
		prev = current
		current = next
	}

	l.head = prev
}

// Delete middle of node in Linked List
func (l *linkedList)DeleteMiddleofNode()  {
	// case1: check empty list or not
	if l.head == nil {
		fmt.Println("Linked list empty!")
		return
	}

	// Get length of Linked list
	len := l.GetLenghtofLinkedList()

	// Get middle value
	mid := len / 2

	// Traverse to middle of node
	current := l.head
	for i := 1; i < mid; i++ {
		current = current.next
	}

	current.next = current.next.next
}
