package rope

type Node struct {
	Left  *Node
	Right *Node
	Size  int
	Data  string
}

func Size(n *Node) int {
	if n == nil {
		return 0 
	}
	return n.Size
	
}