package amocrm

type (
	//RequestParams параметры GET запроса
	AccountRequestParams struct {
		With string
	}
	AccountResponse struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Subdomain      string `json:"subdomain"`
		Currency       string `json:"currency"`
		Timezone       string `json:"timezone"`
		TimezoneOffset string `json:"timezone_offset"`
		Language       string `json:"language"`
		DatePattern    struct {
			Date     string `json:"date"`
			Time     string `json:"time"`
			DateTime string `json:"date_time"`
			TimeFull string `json:"time_full"`
		} `json:"date_pattern"`
		CurrentUser int `json:"current_user"`
		Embedded    struct {
			Users        []AccountUser `json:"users"`
			CustomFields struct {
				Contacts struct {
					Num179391 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
					} `json:"179391"`
					Num179393 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
						Enums struct {
							Num361629 string `json:"361629"`
							Num361631 string `json:"361631"`
							Num361633 string `json:"361633"`
							Num361635 string `json:"361635"`
							Num361637 string `json:"361637"`
							Num361639 string `json:"361639"`
						} `json:"enums"`
					} `json:"179393"`
					Num179395 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
						Enums struct {
							Num361641 string `json:"361641"`
							Num361643 string `json:"361643"`
							Num361645 string `json:"361645"`
						} `json:"enums"`
					} `json:"179395"`
					Num179399 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
						Enums struct {
							Num361647 string `json:"361647"`
							Num361649 string `json:"361649"`
							Num361651 string `json:"361651"`
							Num361653 string `json:"361653"`
							Num361655 string `json:"361655"`
							Num361657 string `json:"361657"`
						} `json:"enums"`
					} `json:"179399"`
					Num474921 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
					} `json:"474921"`
					Num531855 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
					} `json:"531855"`
				} `json:"contacts"`
				Leads struct {
					Num531667 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
						Enums struct {
							Num1054811 string `json:"1054811"`
							Num1054813 string `json:"1054813"`
							Num1054815 string `json:"1054815"`
							Num1054817 string `json:"1054817"`
							Num1054819 string `json:"1054819"`
							Num1054821 string `json:"1054821"`
							Num1054969 string `json:"1054969"`
							Num1054971 string `json:"1054971"`
							Num1060097 string `json:"1060097"`
						} `json:"enums"`
					} `json:"531667"`
				} `json:"leads"`
				Companies struct {
					Num179393 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
						Enums struct {
							Num361629 string `json:"361629"`
							Num361631 string `json:"361631"`
							Num361633 string `json:"361633"`
							Num361635 string `json:"361635"`
							Num361637 string `json:"361637"`
							Num361639 string `json:"361639"`
						} `json:"enums"`
					} `json:"179393"`
					Num179395 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
						Enums struct {
							Num361641 string `json:"361641"`
							Num361643 string `json:"361643"`
							Num361645 string `json:"361645"`
						} `json:"enums"`
					} `json:"179395"`
					Num179397 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
					} `json:"179397"`
					Num179401 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
					} `json:"179401"`
					Num531857 struct {
						ID          int    `json:"id"`
						Name        string `json:"name"`
						FieldType   int    `json:"field_type"`
						Sort        int    `json:"sort"`
						IsMultiple  bool   `json:"is_multiple"`
						IsSystem    bool   `json:"is_system"`
						IsEditable  bool   `json:"is_editable"`
						IsRequired  bool   `json:"is_required"`
						IsDeletable bool   `json:"is_deletable"`
						IsVisible   bool   `json:"is_visible"`
						Params      struct {
						} `json:"params"`
					} `json:"531857"`
				} `json:"companies"`
				Customers []interface{} `json:"customers"`
			} `json:"custom_fields"`
			NoteTypes struct {
				Num1 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"1"`
				Num2 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"2"`
				Num3 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"3"`
				Num4 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"4"`
				Num5 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"5"`
				Num6 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"6"`
				Num7 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"7"`
				Num8 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"8"`
				Num9 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"9"`
				Num10 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"10"`
				Num11 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"11"`
				Num12 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"12"`
				Num13 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"13"`
				Num17 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"17"`
				Num99 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"99"`
				Num101 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"101"`
				Num102 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"102"`
				Num103 struct {
					ID         int    `json:"id"`
					Code       string `json:"code"`
					IsEditable bool   `json:"is_editable"`
				} `json:"103"`
			} `json:"note_types"`
			Groups struct {
				Num0 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"0"`
				Num175591 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"175591"`
				Num180259 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"180259"`
				Num180952 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"180952"`
				Num184612 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"184612"`
			} `json:"groups"`
			TaskTypes struct {
				Num1 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"1"`
				Num2 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"2"`
				Num3 struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"3"`
			} `json:"task_types"`
			Pipelines struct {
				Num1146976 struct {
					ID       int    `json:"id"`
					Name     string `json:"name"`
					Sort     int    `json:"sort"`
					IsMain   bool   `json:"is_main"`
					Statuses struct {
						Num142 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"142"`
						Num143 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"143"`
						Num19743178 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"19743178"`
						Num19743181 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"19743181"`
						Num19743184 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"19743184"`
						Num19743187 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"19743187"`
					} `json:"statuses"`
					Links struct {
						Self struct {
							Href   string `json:"href"`
							Method string `json:"method"`
						} `json:"self"`
					} `json:"_links"`
				} `json:"1146976"`
				Num1209508 struct {
					ID       int    `json:"id"`
					Name     string `json:"name"`
					Sort     int    `json:"sort"`
					IsMain   bool   `json:"is_main"`
					Statuses struct {
						Num142 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"142"`
						Num143 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"143"`
						Num20307982 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20307982"`
						Num20307985 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20307985"`
						Num20307988 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20307988"`
					} `json:"statuses"`
					Links struct {
						Self struct {
							Href   string `json:"href"`
							Method string `json:"method"`
						} `json:"self"`
					} `json:"_links"`
				} `json:"1209508"`
				Num1234921 struct {
					ID       int    `json:"id"`
					Name     string `json:"name"`
					Sort     int    `json:"sort"`
					IsMain   bool   `json:"is_main"`
					Statuses struct {
						Num142 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"142"`
						Num143 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"143"`
						Num20581042 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20581042"`
						Num20581045 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20581045"`
						Num20581669 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20581669"`
						Num20583469 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20583469"`
						Num20634835 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20634835"`
					} `json:"statuses"`
					Links struct {
						Self struct {
							Href   string `json:"href"`
							Method string `json:"method"`
						} `json:"self"`
					} `json:"_links"`
				} `json:"1234921"`
				Num1240990 struct {
					ID       int    `json:"id"`
					Name     string `json:"name"`
					Sort     int    `json:"sort"`
					IsMain   bool   `json:"is_main"`
					Statuses struct {
						Num142 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"142"`
						Num143 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"143"`
						Num20633266 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20633266"`
						Num20633269 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20633269"`
						Num20633272 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20633272"`
					} `json:"statuses"`
					Links struct {
						Self struct {
							Href   string `json:"href"`
							Method string `json:"method"`
						} `json:"self"`
					} `json:"_links"`
				} `json:"1240990"`
				Num1252699 struct {
					ID       int    `json:"id"`
					Name     string `json:"name"`
					Sort     int    `json:"sort"`
					IsMain   bool   `json:"is_main"`
					Statuses struct {
						Num142 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"142"`
						Num143 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"143"`
						Num20738692 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20738692"`
						Num20738695 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20738695"`
						Num20738698 struct {
							ID         int    `json:"id"`
							Name       string `json:"name"`
							Color      string `json:"color"`
							Sort       int    `json:"sort"`
							IsEditable bool   `json:"is_editable"`
						} `json:"20738698"`
					} `json:"statuses"`
					Links struct {
						Self struct {
							Href   string `json:"href"`
							Method string `json:"method"`
						} `json:"self"`
					} `json:"_links"`
				} `json:"1252699"`
			} `json:"pipelines"`
		} `json:"_embedded"`
	}

	AccountUser struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Login    string `json:"login"`
		Language string `json:"language"`
		GroupID  int    `json:"group_id"`
		IsActive bool   `json:"is_active"`
		IsFree   bool   `json:"is_free"`
		IsAdmin  bool   `json:"is_admin"`
		Rights   struct {
			Mail          string `json:"mail"`
			IncomingLeads string `json:"incoming_leads"`
			Catalogs      string `json:"catalogs"`
			LeadAdd       string `json:"lead_add"`
			LeadView      string `json:"lead_view"`
			LeadEdit      string `json:"lead_edit"`
			LeadDelete    string `json:"lead_delete"`
			LeadExport    string `json:"lead_export"`
			ContactAdd    string `json:"contact_add"`
			ContactView   string `json:"contact_view"`
			ContactEdit   string `json:"contact_edit"`
			ContactDelete string `json:"contact_delete"`
			ContactExport string `json:"contact_export"`
			CompanyAdd    string `json:"company_add"`
			CompanyView   string `json:"company_view"`
			CompanyEdit   string `json:"company_edit"`
			CompanyDelete string `json:"company_delete"`
			CompanyExport string `json:"company_export"`
			TaskEdit      string `json:"task_edit"`
			TaskDelete    string `json:"task_delete"`
		} `json:"rights"`
	}
)
