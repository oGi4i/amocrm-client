# AMOCRM Go Client

## Examples
```
import(
    "fmt"
    "github.com/javdet/amocrm-client"
)
//Connect AMO
amo, err := amocrm.New("https://example.amocrm.ru", "example@gmail.com", "453af17f1fdsfsd7792aec4676690567")

// Search contact
contacts, err := amo.GetContact(amocrm.ContactRequestParams{Query: "79999999999"})
fmt.Println(contacts)
fmt.Println(contacts[0].ID)

//Add note
resp, err := amo.AddNote(amocrm.Note{ElementID: contacts[0].ID, ElementType: 1, NoteType: 4, Text: "test4"})

// Get Account info
account, err := amo.GetAccount(amocrm.AccountRequestParams{With: "users"})

//Add call notify
resp, err := amo.AddEvent(amocrm.Event{PhoneNumber: "79999999999", Type: "phone_call", Users: []string{"user_id"}})

//Add lead
resp, err := amo.AddLead(amocrm.Lead{Name: "Call to XXXXX", StatusID: "12345687", ResponsibleUserID: "123456", ContactsID: []string{"24248411"}})

//Add Incoming lead
resp, err := amo.AddIncomingLeadCall(
    amocrm.IncomingLead {
        SourceName: "call from 98234377", 
        SourceUID: "8e64ba2e8822ba378", 
        IncomingEntities: amocrm.IncomingEntities {
            Leads: []amocrm.IncomingLeadParams{
                amocrm.IncomingLeadParams{
                    Name: "call from 98234377",
                },
            },
        }, 
        IncomingLeadInfo: amocrm.IncomingLeadInfo{
            To: "102", 
            From: "73433859994", 
            DateCall: time.Now().Unix(), 
            Duration: "60", 
            Link: "https://callcenter.dela.bz/test1.mp3", 
            ServiceCode: "delabz_widget", 
            Uniq: "8e64ba2e883389",
            AddNote: true,
        },
    },
)

//Get Pipelines
pipelines, err := amo.GetPipelines(amocrm.PipelineRequestParams{Id: ""})
```
