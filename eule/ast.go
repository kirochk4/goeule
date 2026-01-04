package eule

type astNode interface {
	astNodeMark()
}

type astDecl interface {
	astNode
	astDeclMark()
}

type astStmt interface {
	astNode
	astStmtMark()
}

type astExpr interface {
	astNode
	astExprMark()
}

type block = []astDecl

/* == declarations ========================================================== */

type errorDecl struct {
	err error
}

type varDecl = struct {
	name varName
	init astExpr
}
type variableDecl struct {
	vars []varDecl
}

type functionDecl struct {
	name     varName
	function *functionLit
}

type stmtDecl struct {
	stmt astStmt
}

/* == statements ============================================================ */

type blockStmt struct {
	block
}

type emptyStmt empty

type ifStmt struct {
	init  astDecl // Variable declaration or expression.
	cond  astExpr
	then  astStmt
	else_ astStmt
}

type forStmt struct {
	init astDecl // Variable declaration or expression.
	cond astExpr
	post astExpr
	loop astStmt
}

type forEachStmt struct {
	// TODO: implement
}

type whileStmt struct {
	cond astExpr
	loop astStmt
}

type doStmt struct {
	loop astStmt
	cond astExpr
}

type continueStmt empty

type breakStmt empty // TODO: Label jump.

type throwStmt struct {
	throw astExpr
}

type tryStmt struct {
	try     astStmt
	catch   astStmt
	as      varName
	finally astStmt
}

type returnStmt struct {
	value astExpr
}

type exprStmt struct {
	expr astExpr
}

/* == expression ============================================================ */

type emptyExpr empty

type assignExpr struct {
	left  astExpr
	right astExpr
}

type prefixExpr struct {
	op    token
	right astExpr
}

type infixExpr struct {
	left  astExpr
	op    token
	right astExpr
}

type postfixExpr struct {
	left astExpr
	op   token
}

type callExpr struct {
	left astExpr
	args []astExpr
}

// Also used for dot properties.
type indexExpr struct {
	left  astExpr
	index astExpr
}

type protoTableExpr struct {
	proto astExpr
	table *tableLit
}

/* ==literals =============================================================== */

type identifierLit struct {
	varName
}

type nihilLit empty

type booleanLit struct {
	value bool
}

type integerLit struct {
	value int64
}

type floatLit struct {
	value float64
}

type stringLit struct {
	value string
}

type tableLit struct {
	pairs map[astExpr]astExpr
	array []astExpr
}

type functionLit struct {
	params []varName
	body   block
}

/* == marks ================================================================= */

func (n *errorDecl) astDeclMark()    {}
func (n *variableDecl) astDeclMark() {}
func (n *functionDecl) astDeclMark() {}
func (n *stmtDecl) astDeclMark()     {}

func (n *emptyStmt) astStmtMark()    {}
func (n *blockStmt) astStmtMark()    {}
func (n *ifStmt) astStmtMark()       {}
func (n *forStmt) astStmtMark()      {}
func (n *forEachStmt) astStmtMark()  {}
func (n *whileStmt) astStmtMark()    {}
func (n *doStmt) astStmtMark()       {}
func (n *continueStmt) astStmtMark() {}
func (n *breakStmt) astStmtMark()    {}
func (n *throwStmt) astStmtMark()    {}
func (n *tryStmt) astStmtMark()      {}
func (n *returnStmt) astStmtMark()   {}
func (n *exprStmt) astStmtMark()     {}

func (n *emptyExpr) astExprMark()      {}
func (n *assignExpr) astExprMark()     {}
func (n *prefixExpr) astExprMark()     {}
func (n *infixExpr) astExprMark()      {}
func (n *postfixExpr) astExprMark()    {}
func (n *callExpr) astExprMark()       {}
func (n *indexExpr) astExprMark()      {}
func (n *protoTableExpr) astExprMark() {}
func (n *identifierLit) astExprMark()  {}
func (n *nihilLit) astExprMark()       {}
func (n *booleanLit) astExprMark()     {}
func (n *integerLit) astExprMark()     {}
func (n *floatLit) astExprMark()       {}
func (n *stringLit) astExprMark()      {}
func (n *tableLit) astExprMark()       {}
func (n *functionLit) astExprMark()    {}

// ==  ==  ==  ==

func (n *errorDecl) astNodeMark()    {}
func (n *variableDecl) astNodeMark() {}
func (n *functionDecl) astNodeMark() {}
func (n *stmtDecl) astNodeMark()     {}

func (n *emptyStmt) astNodeMark()    {}
func (n *blockStmt) astNodeMark()    {}
func (n *ifStmt) astNodeMark()       {}
func (n *forStmt) astNodeMark()      {}
func (n *forEachStmt) astNodeMark()  {}
func (n *whileStmt) astNodeMark()    {}
func (n *doStmt) astNodeMark()       {}
func (n *continueStmt) astNodeMark() {}
func (n *breakStmt) astNodeMark()    {}
func (n *throwStmt) astNodeMark()    {}
func (n *tryStmt) astNodeMark()      {}
func (n *returnStmt) astNodeMark()   {}
func (n *exprStmt) astNodeMark()     {}

func (n *emptyExpr) astNodeMark()      {}
func (n *assignExpr) astNodeMark()     {}
func (n *prefixExpr) astNodeMark()     {}
func (n *infixExpr) astNodeMark()      {}
func (n *postfixExpr) astNodeMark()    {}
func (n *callExpr) astNodeMark()       {}
func (n *indexExpr) astNodeMark()      {}
func (n *protoTableExpr) astNodeMark() {}
func (n *identifierLit) astNodeMark()  {}
func (n *nihilLit) astNodeMark()       {}
func (n *booleanLit) astNodeMark()     {}
func (n *integerLit) astNodeMark()     {}
func (n *floatLit) astNodeMark()       {}
func (n *stringLit) astNodeMark()      {}
func (n *tableLit) astNodeMark()       {}
func (n *functionLit) astNodeMark()    {}

/* == printer =============================================================== */

type printer struct{}

func (p *printer) print([]astDecl) {}
