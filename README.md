# AMOCRM Go Client

## Examples
```
import(
    "fmt"
    "github.com/javdet/amocrm-client"
)
amo, err := amocrm.New("https://example.amocrm.ru", "example@gmail.com", "453af17f1fdsfsd7792aec4676690567")
contacts, err := amo.GetContact(amocrm.ContactRequestParams{Query: "79999999999"})
fmt.Println(contacts)
fmt.Println(contacts[0].ID)

resp, err := amo.AddNote(amocrm.Note{ElementID: contacts[0].ID, ElementType: 1, NoteType: 4, Text: "test4"})

account, err := amo.GetAccount(amocrm.AccountRequestParams{With: "users"})
```