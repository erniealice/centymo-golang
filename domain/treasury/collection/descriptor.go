package collection

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.collection",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "treasury_collection"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "collection.json", Key: "collection"},
		LabelName: "CollectionLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "collection:list",
			AppEntry: &compose.AppEntry{
				Key: "cash", Route: "collection.list", Params: map[string]string{"status": "pending"},
				Label: "Cash", Icon: "icon-credit-card",
				Permission: "collection:list",
			},
			Items: []compose.NavItem{
				// cash app — Collections section
				{Key: "collections-pending", Route: "collection.list", Params: map[string]string{"status": "pending"},
					Label: "Pending", Icon: "icon-clock", Permission: "collection:list"},
				{Key: "collections-completed", Route: "collection.list", Params: map[string]string{"status": "completed"},
					Label: "Complete", Icon: "icon-check-circle", Permission: "collection:list"},
			},
		},
	}
}
