package tree

type BranchTree struct {
	Area            string `json:"area"`
	BranchName      string `json:"branchName"`
	BranchID        string `json:"branchId"`
	PatentName      string `json:"patentName"`
	PatentID        string `json:"patentId"`
	CurrectBranchId string `json:"currectBranchId"`
}

type NodeExample1 struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	EmployeeBelongs bool            `json:"employeeBelongs"`
	Children        []*NodeExample1 `json:"child,omitempty"`
}

// var (
// 	nodeTable = map[string]*NodeExample1{}
// 	root      NodeExample1
// )

// type Node struct {
// 	tag      string
// 	id       string
// 	class    string
// 	children []*Node
// }

type NodeInterface interface {
	AssembleTree(BranchTree []BranchTree) error
}

func AssembleTreeHandler(BranchTree []BranchTree) interface{} {

	var Node NodeInterface

	Node = &NodeExample1{Name: "Root", Children: []*NodeExample1{}}
	Node.AssembleTree(BranchTree)
	////////_, Result := Node.AssembleTree(BranchTree)

	Node = &NodeExample2{"", "", "root", false, 0, nil}
	Node.AssembleTree(BranchTree)
	///////Node, _ := AssembleTreeExample2(BranchTree)

	///////Node = findByIdDFSInterFace(&NodeExample2{"", "", "root", false, 0, nil}, "187ac25d-ef9d-11eb-9114-005056a2ef46")
	///////Node = findByIdDFSInterFace(Node, "187ac25d-ef9d-11eb-9114-005056a2ef46")

	return Node
}

func (node *NodeExample1) AssembleTree(BranchTree []BranchTree) error {

	// root = NodeExample1{Name: "Root", Children: []*NodeExample1{}}
	// nodeTable[""] = &root

	nodeTable := make(map[string]*NodeExample1)
	nodeTable[""] = node

	CurrectBranchId := ""

	for _, value := range BranchTree {
		node.Add(value.BranchID, value.BranchName, value.PatentID, nodeTable)
		CurrectBranchId = value.CurrectBranchId
	}

	// fmt.Printf("main: reading input from stdin\n")
	// scan()

	node.findByIdDFS(CurrectBranchId)
	//node.findById(CurrectBranchId)
	// rootFind := findByIdDFS(node, CurrectBranchId)
	// //rootFind := findById(&root, "187ac25d-ef9d-11eb-9114-005056a2ef46")
	// if rootFind != nil {
	// 	rootFind.EmployeeBelongs = true
	// }

	//fmt.Printf("main: reading input from stdin -- done\n")
	//showExample1()
	//fmt.Printf("main: end\n")

	// byteTest, err := json.Marshal(root)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(byteTest))

	//return nil, root
	return nil
}

func (node *NodeExample1) Add(id, name, parentId string, nodeTable map[string]*NodeExample1) {
	//fmt.Printf("add: id=%v name=%v parentId=%v\n", id, name, parentId)

	node = &NodeExample1{Id: id, Name: name, Children: []*NodeExample1{}}

	if parentId == "mock" {
		//root = node
	} else {

		parent, ok := nodeTable[parentId]
		if !ok {
			//fmt.Printf("add: parentId=%v: not found\n", parentId)
			return
		}

		parent.Children = append(parent.Children, node)
	}

	nodeTable[id] = node
}

func (node *NodeExample1) findByIdDFS(id string) {
	if node.Id == id {
		node.EmployeeBelongs = true
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			child.findByIdDFS(id)
		}
	}
}

func (node *NodeExample1) findById(id string) {
	queue := make([]*NodeExample1, 0)
	queue = append(queue, node)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.Id == id {
			nextUp.EmployeeBelongs = true
		}
		if len(nextUp.Children) > 0 {
			for _, child := range nextUp.Children {
				queue = append(queue, child)
			}
		}
	}

}

// func scanExample1() {
// 	input := os.Stdin
// 	reader := bufio.NewReader(input)
// 	lineCount := 0
// 	for {
// 		lineCount++
// 		line, err := reader.ReadString('\n')
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			fmt.Printf("error reading lines: %v\n", err)
// 			return
// 		}
// 		tokens := strings.Fields(line)
// 		if t := len(tokens); t != 3 {
// 			fmt.Printf("bad input line %v: tokens=%d [%v]\n", lineCount, t, line)
// 			continue
// 		}
// 		addExample1(tokens[0], tokens[1], tokens[2])
// 	}
// }

// func showNodeExample1(node *NodeExample1, prefix string) {
// 	if prefix == "" {
// 		fmt.Printf("%v -> %v -> %v \n\n", node.Name, node.Id, node.EmployeeBelongs)
// 	} else {
// 		fmt.Printf("%v %v -> %v -> %v \n\n", prefix, node.Name, node.Id, node.EmployeeBelongs)
// 	}
// 	for _, n := range node.Children {
// 		showNodeExample1(n, prefix+"--")
// 	}
// }

// func showExample1() {
// 	if &root == nil {
// 		fmt.Printf("show: root node not found\n")
// 		return
// 	}
// 	fmt.Printf("RESULT:\n")
// 	showNodeExample1(&root, "")
// }

// TODO: Сделать 2 варианта через интрефейс переменную и через просто пустой интерфейс параметр
func findById(root *NodeExample1, id string) *NodeExample1 {
	queue := make([]*NodeExample1, 0)
	queue = append(queue, root)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.Id == id {
			return nextUp
		}
		if len(nextUp.Children) > 0 {
			for _, child := range nextUp.Children {
				queue = append(queue, child)
			}
		}
	}
	return nil
}

// func findByIdDFSProblems(node *NodeExample1, id string) *NodeExample1 {
// 	if node.Id == id {
// 		return node
// 	}

// 	if len(node.Children) > 0 {
// 		for _, child := range node.Children {
// 			findByIdDFSProblems(child, id)
// 		}
// 	}
// 	return nil
// }

func findByIdDFS(node *NodeExample1, id string) *NodeExample1 {
	if node.Id == id {
		return node
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			NodeExampleInside := findByIdDFS(child, id)
			if NodeExampleInside != nil {
				return NodeExampleInside
			}
		}
	}
	return nil
}

func findByIdDFSInterFace(nodeArg interface{}, id string) *NodeExample2 {

	node, ok := nodeArg.(*NodeExample2)
	if !ok {
		return nil
	}

	if node.Id == id {
		return node
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			NodeExampleInside := findByIdDFSInterFace(child, id)
			if NodeExampleInside != nil {
				return NodeExampleInside
			}
		}
	}
	return nil
}

//--------------------------------------------------------------------------------------------------------------------------------------------------------

type NodeExample2 struct {
	Id              string          `json:"id"`
	ParentId        string          `json:"-"`
	Name            string          `json:"name"`
	EmployeeBelongs bool            `json:"employeeBelongs"`
	Leaf            int             `json:"leaf"` //`json:"leaf,omitempty"`
	Children        []*NodeExample2 `json:"child,omitempty"`
}

func (this *NodeExample2) Size() int {
	var size int = len(this.Children)
	for _, c := range this.Children {
		size += c.Size()
	}
	return size
}
func (this *NodeExample2) Add(nodes ...*NodeExample2) bool {
	var size = this.Size()
	for _, n := range nodes {
		if n.ParentId == this.Id {
			this.Children = append(this.Children, n)
		} else {
			for _, c := range this.Children {
				if c.Add(n) {
					break
				}
			}
		}
	}
	this.Leaf = this.Size()
	return this.Size() == size+len(nodes)
}

func (node *NodeExample2) AssembleTree(BranchTree []BranchTree) error {

	CurrectBranchId := ""
	var data []*NodeExample2
	for _, value := range BranchTree {
		data = append(data, &NodeExample2{Id: value.BranchID, ParentId: value.PatentID, Name: value.BranchName})
		//CurrectBranchId = value.CurrectBranchId
	}

	node.Add(data...)

	node.findByIdDFS(CurrectBranchId)
	//node.findById(CurrectBranchId)

	// fmt.Println(node.Add(data...), node.Size())
	// bytes, _ := json.MarshalIndent(node, "", "\t") //formated output
	// //bytes, _ := json.Marshal(root)
	// fmt.Println(string(bytes))

	// rootFind := findByIdDFS(&root, CurrectBranchId)
	// //rootFind := findById(&root, "187ac25d-ef9d-11eb-9114-005056a2ef46")
	// if rootFind != nil {
	// 	rootFind.EmployeeBelongs = true
	// }

	return nil
}

func (node *NodeExample2) findByIdDFS(id string) {
	if node.Id == id {
		node.EmployeeBelongs = true
	}

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			child.findByIdDFS(id)
		}
	}
}

func (node *NodeExample2) findById(id string) {
	queue := make([]*NodeExample2, 0)
	queue = append(queue, node)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.Id == id {
			nextUp.EmployeeBelongs = true
		}
		if len(nextUp.Children) > 0 {
			for _, child := range nextUp.Children {
				queue = append(queue, child)
			}
		}
	}

}

// func AssembleTreeExample2(BranchTree []store.BranchTree) (NodeExample2, error) {

// 	//var root *NodeExample2 = &NodeExample2{"", "", "root", "", nil}

// 	root := NodeExample2{"", "", "root", false, "", nil}

// 	var data []*NodeExample2
// 	for _, value := range BranchTree {
// 		data = append(data, &NodeExample2{Id: value.BranchID, ParentId: value.PatentID, Name: value.BranchName})
// 	}

// 	// data := []*NodeExampl2{
// 	// 	&NodeExampl2{"002", "001", "Shooping", "0", nil},
// 	// 	&NodeExampl2{"003", "002", "Housewares", "0", nil},
// 	// 	&NodeExampl2{"004", "003", "Kitchen", "1", nil},
// 	// 	&NodeExampl2{"005", "003", "Officer", "1", nil},
// 	// 	&NodeExampl2{"006", "002", "Remodeling", "0", nil},
// 	// 	&NodeExampl2{"007", "006", "Retile kitchen", "1", nil},
// 	// 	&NodeExampl2{"008", "006", "Paint bedroom", "1", nil},
// 	// 	&NodeExampl2{"009", "008", "Ceiling", "1", nil},
// 	// 	&NodeExampl2{"010", "006", "Other", "1", nil},
// 	// 	&NodeExampl2{"011", "001", "Misc", "1", nil},
// 	// }
// 	fmt.Println(root.Add(data...), root.Size())
// 	bytes, _ := json.MarshalIndent(root, "", "\t") //formated output
// 	//bytes, _ := json.Marshal(root)
// 	fmt.Println(string(bytes))

// 	return root, nil
// }
