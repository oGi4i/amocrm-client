# AmoCRM Go Client

## Examples
```
import(
    "fmt"

    "github.com/ogi4i/amocrm-client"
)
//Connect AMO
amo, err := amocrm.New("https://example.amocrm.ru", "example@gmail.com", "453af17f1fdsfsd7792aec4676690567")

// Search contact
contacts, err := amo.GetContacts(amocrm.ContactRequestParams{Query: "79999999999"})
fmt.Println(contacts)
fmt.Println(contacts[0].ID)

//Add note
resp, err := amo.AddNote(amocrm.Note{ElementID: contacts[0].ID, ElementType: 1, NoteType: 4, Text: "test4"})

// Get Account info
account, err := amo.GetAccount(amocrm.AccountRequestParams{With: "users"})

//Add lead
resp, err := amo.AddLead(amocrm.Lead{Name: "Call to XXXXX", StatusID: "12345687", ResponsibleUserID: "123456", ContactsID: []string{"24248411"}})

//Get Pipelines
pipelines, err := amo.GetPipelines(amocrm.PipelineRequestParams{ID: ""})
```
