package dao
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "fmt"
    "os"
    "errors"
) 
type Model struct {
    Data string
    Token string
    Pass1 string
}
func Create() {
   fmt.Println("creating database");
   os.Remove("./data.db")

   db := _OpenDb()
   if db == nil {
       log.Fatal("Could not open db")
   }
   defer db.Close()
   sql := `
    create table data (id integer not null primary key, data text, token text, pass_phrase1);
   delete from data;
    `
    _, err := db.Exec(sql)
    if err != nil {
        log.Printf("%q: %s\n", err, sql)
        fmt.Printf("%q: %s\n", err, sql)
        return
    }
}
func _OpenDb() (*sql.DB) {
    db, err := sql.Open("sqlite3", "./data.db")
    if err != nil {
       log.Fatal(err)
       fmt.Println(err)
       return nil
    }
    return db

}
func Insert(model Model) (bool, error){
    db := _OpenDb()
    if db == nil {
       log.Fatal("Could not open db")
       fmt.Printf("Could not open db");
    }
    sql := "insert into data (data, token, pass_phrase1) values (?, ?, ?)"
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
        return false, errors.New(fmt.Sprint("Could not begin transaction ", err.Error))
    }
    stmt, err := tx.Prepare(sql)
    if err != nil {
        fmt.Printf("prepare:%v\n", err)
        return false, err
        return false, errors.New(fmt.Sprint("Could not prepare statement ", err.Error))
    }
    defer stmt.Close()
    _, err = stmt.Exec(model.Data, model.Token, model.Pass1)
    if err != nil {
        fmt.Printf("exec: %v\n", err)
        return false, errors.New(fmt.Sprint("Could not execute statement ", err.Error))
    }
    tx.Commit()
    return true, nil
}

func Select(token string) (*Model){
    db := _OpenDb()
    if db == nil {
       fmt.Println("Could not open db")
       return nil
    }
    stmt, err := db.Prepare("select data, token, pass_phrase1 from data where token = ?")
    if err != nil {
       fmt.Println("Could not prepare statement")
    }
    defer stmt.Close()
    model := &Model{}
    err = stmt.QueryRow(token).Scan(&model.Data, &model.Token, &model.Pass1)
    if err != nil {
        log.Fatal(err)
    }
    db.Close()
    return model
}

func Delete(token string) (bool) {
    db := _OpenDb()
    if db == nil {
       fmt.Println("Could not open db")
       return false
    }
    stmt, err := db.Prepare("delete from data where token = ?")
    if err != nil {
       fmt.Println("Could not prepare statement")
    }
    defer stmt.Close()
    _, err = stmt.Exec(token)
    if err != nil {
        fmt.Printf("exec: %v\n", err)
        return false
    }
    if err != nil {
        log.Fatal(err)
    }
    db.Close()
    return true
}
