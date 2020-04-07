package main

import (
	"fmt"
	"log"
	"bytes"
	"regexp"
	"net/http"
	"io/ioutil"
	"database/sql"
	"encoding/json"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

func main () {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/translate", receiveWord).Methods("POST")
	router.HandleFunc("/translate/{word}", translatedWord).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}


func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Oportunidade Mercado Pago - Julio Verne")
}


type Translation struct {
	Words string `json:"translation"`
}

func receiveWord(w http.ResponseWriter, r *http.Request) {
	db := connect()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid content")
	}

	var translate map[string]interface {}
	json.Unmarshal([]byte(body), &translate)

	words := translate["words"].([]interface {})
	count_words := len(words)

	var letterArray bytes.Buffer

	for i := 0; i < count_words; i++ {
		word := words[i].(map[string]interface {})
		letters := word["letters"].([]interface {})
		var wordsToDb bytes.Buffer
		
		//insert
		insertWord, err := db.Query("INSERT INTO words (word) VALUES (null)")
		if err != nil {
            panic(err.Error())
		}
		defer insertWord.Close()

		// select
		var lastId DbWord
		selectLastId := 0
		selectId := db.QueryRow("SELECT id FROM words ORDER BY id DESC LIMIT 1")
		
		switch errr := selectId.Scan(&lastId.ID); errr {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			selectLastId = lastId.ID
		default:
			panic(err.Error())
	  	}

		for j := 0; j < len(letters); j++ {
			letter := translateLetter(letters[j].(string))
			letterArray.WriteString(letter)
			wordsToDb.WriteString(letter)
			
			insertLetter, err := db.Query(
				"INSERT INTO letters (letter_it, letter_translated, word_id) VALUES (?, ?, ?)",
				letters[j].(string),
				letter,
				selectLastId,
			)
			if err != nil {
				panic(err.Error())
			}
			defer insertLetter.Close()
		}
		
		// update
		updateWord, errrr := db.Query("UPDATE words SET word = ? WHERE id = ?", wordsToDb.String(), selectLastId)
		if errrr != nil {
            panic(errrr.Error())
		}
		defer updateWord.Close()

		// space between words
		if i != (count_words - 1) {
			letterArray.WriteString(" ")
		}
	}

	translation := &Translation{
		Words: letterArray.String(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(translation)
}


type Letters struct {
	Letters []string `json:"letters"`
}

type Word struct {
	Word Letters `json:"word"`
}

type Untranslate struct {
	KnownWord bool  `json:"knownWord"`
	Versions []Word `json:"versions"`
}

func translatedWord(w http.ResponseWriter, r *http.Request) {
	db := connect()
	words := []Word{}
	knownWord := false
	request_word := mux.Vars(r)["word"]

	select_words, err := db.Query("SELECT id, word FROM words WHERE word = ?", request_word)
	if err != nil {
		panic(err.Error())
	}

	for select_words.Next() {
		var dbwords DbWord
		err := select_words.Scan(&dbwords.ID, &dbwords.Word)
		if err != nil {
			panic(err.Error())
		}

		// select letters
		select_letters, errr := db.Query("SELECT letter_it FROM letters WHERE word_id = ?", dbwords.ID)
		if errr != nil {
			panic(err.Error())
		}

		letters := []string{}
		for select_letters.Next() {
			var dbletters DbLetter
			errrr := select_letters.Scan(&dbletters.Letter_it)
			if errrr != nil {
				panic(errrr.Error())
			}
			letters = append(letters, dbletters.Letter_it)	
		}

		letter := Letters{
			Letters: letters,
		}

		word := Word{
			Word: letter,
		}
		words = append(words, word)

		knownWord = true
	}

	response := Untranslate{
		KnownWord: knownWord,
		Versions: words,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}


func translateLetter(letter string) string {
	vogals := regexp.MustCompile(`a|e|i|o|u`)
	consonants := regexp.MustCompile(`b|c|d|f|g|h|j|k|l|m|n|p|q|r|s|t|v|w|x|y|z`)

	var count_vogals []string = vogals.FindAllString(letter, -1)
	var count_consonants []string = consonants.FindAllString(letter, -1)
	
	ascii := len(count_vogals) - len(count_consonants)
	result := string(ascii)

	return result
}


type DbWord struct {
	ID int
	Word string
}

type DbLetter struct {
	ID int
	Letter_it string
	Letter_translated string
	Word_id int
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/meli_translation")
	if err != nil {
        panic(err.Error())
	}
	
	return db
}