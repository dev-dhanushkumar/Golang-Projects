package main

import (
	"github.com/dev-dhanushkumar/Golang-Projects/GOLANG/Linked_List_Impl/Linkedlist"
	"fmt"
)

func main()  {
	list := Linkedlist.New()	

	list.InsertAtBeginning(10)
	list.InsertAtBeginning(20)
	list.InsertAtEnd(30)
	list.InsertAtBeginning(40)
	list.InsertAtMiddle(50)
	list.InsertAtIdx(100,9)

	list.Display()
	fmt.Println("Lenght of Linked List: ", list.GetLenghtofLinkedList())


	// Reverse the LinkedList
	list.ReverseList()
	list.Display()

	// Remove middle of Node in list
	list.DeleteMiddleofNode()
	list.Display()

	// Delete End Node
	list.DeleteEndofNode()
	list.Display()

	// Delete beginning Node
	list.DeleteBeginningNode()
	list.Display()
}
