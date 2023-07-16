package main

/*
 =========================================
|                                         |
|Queue Implementation for request handling|
|                                         |
 =========================================
*/

type Node struct {
	packets []Packet
	next    *Node
}

type Queue struct {
	head *Node
	tail *Node
	size int
}

func (qe *Queue) AddPack(packets []Packet) {
	newNode := &Node{packets: packets}
	if qe.size == 0 {
		qe.head = newNode
		qe.tail = newNode
	} else {
		qe.tail.next = newNode
		qe.tail = newNode
	}
	qe.size++
}

func (qe *Queue) SendPack() ([]Packet, bool) {
	if qe.size == 0 {
		return nil, false
	}
	packets := qe.head.packets
	qe.head = qe.head.next
	qe.size--
	return packets, true
}

func (qe *Queue) Size() int {
	return qe.size
}
