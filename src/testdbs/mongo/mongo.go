package mongo

import (
	"strconv"

	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"strings"
)

type Person struct {
	Name string
}

func InsertData(address, username, password, database string) error {
	var link string
	link = "mongodb://" + username + ":" + password + "@" + address + "/" + database
	session, err := mgo.Dial(link)
	if err != nil {
		return err
	}
	if session == nil {
		err = fmt.Errorf("InsertData: Dial: session is nil")
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(database).C("testCo")
	for i := 1; i <= 1000; i++ {
		err = c.Insert(&Person{"number" + strconv.Itoa(i)})
		if err != nil {
			return err
		}
	}

	log.Println("InsertData: Insert datas completely")
	return nil
}

func FindPrimary(address, username, password, database string) error {

	addresses := strings.Split(address, ",")

	for _, ad := range addresses {
		link := "mongodb://" + username + ":" + password + "@" + ad + "/" + database
		session, err := mgo.Dial(link)
		if err != nil {
			log.Printf("FindPrimary: %s", err)
			continue
		} else {
			log.Printf("FindPrimary: Now the primary point's address is: [ %s ]\n", ad)
			session.Close()
			return nil
		}
		if session == nil {
			log.Printf("FindPrimary: Dial: [ %s ] session is nil\n", ad)
			continue
		}
	}

	err := fmt.Errorf("FindPrimary: Please keep one mongo point on at least")
	return err
}

func FindData(address, username, password, database string) error {

	var link string
	var people []Person

	link = "mongodb://" + username + ":" + password + "@" + address + "/" + database
	session, err := mgo.Dial(link)

	if err != nil {
		return err
	}
	if session == nil {
		err = fmt.Errorf("FindData: Dial: session is nil")
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(database).C("testCo")

	err = c.Find(nil).All(&people)
	if err != nil {
		log.Println("FindData: Failed to find data")
		return err
	} else {
		log.Printf("FindData: find [ %d ] data successfully\n", len(people))
	}
	return nil
}
