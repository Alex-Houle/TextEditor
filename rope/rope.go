package rope

type Node struct {
	Left  *Node
	Right *Node
	Size  int
	Data  string
}