package main

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

const (
	userRegisterUsernameCheckRule = `
rule UserRegisterUsernameCheck {
    when
      User.Name != "" && User.Password != ""
	then
	  Repo.CreateUser(User);
}
`
)

// User User
type User struct {
	Name     string
	Password string
}

type Repo struct {
	m []*User
}

func (rp *Repo) CreateUser(u *User) {
	rp.m = append(rp.m, u)
}

func main() {
	user := &User{
		Name:     "name",
		Password: "password",
	}

	dataContext := ast.NewDataContext()
	err := dataContext.Add("User", user)
	if err != nil {
		panic(err)
	}

	rp := &Repo{}
	err = dataContext.Add("Repo", rp)
	if err != nil {
		panic(err)
	}

	lib := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(lib)

	err = ruleBuilder.BuildRuleFromResource("demo", "v1", pkg.NewBytesResource([]byte(userRegisterUsernameCheckRule)))
	if err != nil {
		panic(err)
	}

	kb := lib.NewKnowledgeBaseInstance("demo", "v1")
	eng1 := &engine.GruleEngine{MaxCycle: 1}
	err = eng1.Execute(dataContext, kb)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(rp.m))
}
