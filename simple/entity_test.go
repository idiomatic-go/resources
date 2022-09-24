package util

import "fmt"

type versionedStruct struct {
	vers  string
	count int
}

var getVersion = func(a any) string {
	e, ok := getEntity(a)
	if !ok {
		return ""
	}
	return e.vers
}

var getEntity = func(a any) (versionedStruct, bool) {
	if a == nil {
		return versionedStruct{}, false
	}
	if data, ok := a.(versionedStruct); ok {
		return data, ok
	}
	return versionedStruct{}, false
}

func ExampleNoEntity() {
	e := CreateVersionedEntity(nil, getVersion)
	fmt.Printf("IsNewVersion [empty] : %v\n", e.IsNewVersion(""))
	fmt.Printf("IsNewVersion [123] : %v\n", e.IsNewVersion("123"))
	fmt.Printf("IsNewVersion [1.2.3] : %v\n", e.IsNewVersion("1.2.3"))
	if e.Get() == nil {
		fmt.Println("Get : valid")
	} else {
		fmt.Println("failure")
	}

	//Output:
	// IsNewVersion [empty] : false
	// IsNewVersion [123] : false
	// IsNewVersion [1.2.3] : false
	// Get : valid
}

func ExampleNoGetFn() {
	s := versionedStruct{vers: "1.2.3", count: 1}
	e := CreateVersionedEntity(&s, nil)
	fmt.Printf("IsNewVersion [empty] : %v\n", e.IsNewVersion(""))
	fmt.Printf("IsNewVersion [123] : %v\n", e.IsNewVersion("123"))
	fmt.Printf("IsNewVersion [1.2.3] : %v\n", e.IsNewVersion("1.2.3"))

	//Output:
	// IsNewVersion [empty] : false
	// IsNewVersion [123] : true
	// IsNewVersion [1.2.3] : true
}

func ExampleValidEntity() {
	s := versionedStruct{vers: "1.2.3", count: 1}
	e := CreateVersionedEntity(&s, getVersion)
	fmt.Printf("IsNewVersion [empty] : %v\n", e.IsNewVersion(""))
	fmt.Printf("IsNewVersion [123] : %v\n", e.IsNewVersion("123"))
	fmt.Printf("IsNewVersion [1.2.3] : %v\n", e.IsNewVersion("1.2.3"))
	fmt.Printf("Version : %v\n", e.GetVersion())
	entity, ok := getEntity(e.Get())
	if ok {
		fmt.Printf("Entity : [%v]\n", entity)
	} else {
		fmt.Println("Entity : invalid")
	}

	//Output:
	// IsNewVersion [empty] : true
	// IsNewVersion [123] : true
	// IsNewVersion [1.2.3] : false
	// Version : 1.2.3
	// Entity : [{1.2.3 1}]
}

func ExampleChangeEntity() {
	s := versionedStruct{vers: "1.2.3", count: 1}
	e := CreateVersionedEntity(&s, getVersion)
	fmt.Printf("Version : %v\n", e.GetVersion())
	fmt.Printf("IsNewVersion [1.2.3] : %v\n", e.IsNewVersion("1.2.3"))
	entity, ok := getEntity(e.Get())
	if ok {
		fmt.Printf("Entity : [%v]\n", entity)
	} else {
		fmt.Println("Entity : invalid")
	}
	e.Set(versionedStruct{vers: "4.5.6", count: 200})
	fmt.Printf("Version : %v\n", e.GetVersion())
	fmt.Printf("IsNewVersion [1.2.3] : %v\n", e.IsNewVersion("1.2.3"))
	entity, ok = getEntity(e.Get())
	if ok {
		fmt.Printf("Entity : [%v]\n", entity)
	} else {
		fmt.Println("Entity : invalid")
	}

	//Output:
	// Version : 1.2.3
	// IsNewVersion [1.2.3] : false
	// Entity : [{1.2.3 1}]
	// Version : 4.5.6
	// IsNewVersion [1.2.3] : true
	// Entity : [{4.5.6 200}]
}
