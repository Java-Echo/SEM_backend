package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type User struct {
	ID          string `json:"ID"`
	Pwd         string `json:"pwd"`
	Information string `json:"information"`
	Identity    string `json:"identity""`
}
type Data struct {
	ID      int32  `json:"ID"`
	Name    string `json:"name"`
	Explain string `json:"explain"`
}
type Review struct {
	ID       int32  `json:"ID"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Revision string `json:"revision"`
}

func main() {
	//实例化echo对象。
	e := echo.New()

	db, err := sql.Open("mysql", "root:200110@tcp(121.40.54.219:3306)/Search System")
	fmt.Println(err)
	// Register account
	e.POST("/register", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ID := c.FormValue("ID")
		Pwd := c.FormValue("pwd")
		Information := c.FormValue("information")
		Identity := c.FormValue("identity")
		// id is primary key
		rows, err := db.Query("select ID from user where ID=?", ID)
		fmt.Println(err)
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err == nil { //primary key has been used
				ret := User{
					ID: "null",
				}
				return c.JSON(http.StatusOK, ret)
			}
		}
		stmt, err := db.Prepare("insert user set ID=?, pwd=?, information=?, identity=?")
		res, err := stmt.Exec(ID, Pwd, Information, Identity)
		ret, err := res.LastInsertId()
		return c.JSON(http.StatusOK, ret)
	})
	// Login
	e.POST("/login", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ID := c.FormValue("ID")
		Pwd := c.FormValue("pwd")
		rows, err := db.Query("select pwd from user where ID=?", ID)
		fmt.Println(err)
		for rows.Next() {
			var passwd string
			err := rows.Scan(&passwd)
			if err == nil { // founded
				// wrong passwd
				if passwd != Pwd {
					ret := User{
						ID:  ID,
						Pwd: "null",
					}
					return c.JSON(http.StatusOK, ret)
				} else {
					ret := User{
						ID:  ID,
						Pwd: Pwd,
					}
					return c.JSON(http.StatusOK, ret)
				}
			}
		}
		// not founded
		ret := User{
			ID:  "null",
			Pwd: "null",
		}
		return c.JSON(http.StatusOK, ret)
	})
	// alter personal information
	e.POST("/Person/alter", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ID := c.FormValue("ID")
		Pwd := c.FormValue("pwd")
		Info := c.FormValue("information")
		stmt, err := db.Prepare("update user set pwd=?, information=? where ID =?")
		fmt.Println(err)
		res, err := stmt.Exec(Pwd, Info, ID)
		ret, err := res.LastInsertId()
		return c.JSON(http.StatusOK, ret)
	})
	// revision manager
	e.POST("/Manager/Review", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ID := c.FormValue("ID")
		Type := c.FormValue("type")
		Name := c.FormValue("name")
		Revision := c.FormValue("revision")
		if Type == "alter" { // alter
			rows, err := db.Query("select ID from user where ID=?", ID)
			fmt.Println(err)
			for rows.Next() {
				var tmp int32
				if err := rows.Scan(&tmp); err != nil { // not founded
					ret := Review{
						Name: "null",
					}
					return c.JSON(http.StatusOK, ret)
				}
			}
			stmt, err := db.Prepare("update data set  name=?, explain=? where ID=? ")
			res, err := stmt.Exec(Name, Revision, ID)
			ret, err := res.LastInsertId()
			return c.JSON(http.StatusOK, ret)

		} else if Type == "insert" { //insert
			rows, err := db.Query("select ID from user where ID=?", ID)
			fmt.Println(err)
			for rows.Next() {
				var tmp int32
				if err := rows.Scan(&tmp); err == nil { // primary key has been used
					ret := Review{
						Name: "null",
					}
					return c.JSON(http.StatusOK, ret)
				}
			}
			stmt, err := db.Prepare("insert data set ID=?, name=?, explain=?")
			res, err := stmt.Exec(ID, Name, Revision)
			ret, err := res.LastInsertId()
			return c.JSON(http.StatusOK, ret)
		} else {
			ret := Data{
				Name: "null",
			}
			return c.JSON(http.StatusOK, ret)
		}
	})

	//search data
	e.POST("/search", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		Name := c.FormValue("name")
		rows, err := db.Query("select * from data where name=?", Name)
		fmt.Println(err)
		for rows.Next() {
			var Name, Explain string
			var ID int32
			if err := rows.Scan(&ID, &Name, &Explain); err == nil { // founded
				data := Data{
					ID:      ID,
					Name:    Name,
					Explain: Explain,
				}
				return c.JSON(http.StatusOK, data)
			}
		}
		// not founded
		data := Data{
			Name: "null",
		}
		return c.JSON(http.StatusOK, data)
	})
	// start
	e.Start(":8080")
}
