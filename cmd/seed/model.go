package main

// ListModel basic list model
type ListModel struct {
	Title  string
	URI    string
	Fields []FieldModel
}

// FieldModel basic field model
type FieldModel struct {
	Name      string
	SchemaXML string
}

// Lists model definition
var listsModel = []ListModel{
	{
		Title:  "Managers",
		Fields: []FieldModel{},
	},
	{
		Title: "Customers",
		Fields: []FieldModel{
			{
				Name:      "RoutingNumber",
				SchemaXML: `<Field Type="Text" DisplayName="Routing Number" Required="FALSE" MaxLength="9" />`,
			},
			{
				Name: "Region",
				SchemaXML: `<Field Type="Choice" DisplayName="Region" Required="FALSE" Format="Dropdown" FillInChoice="FALSE">
					<CHOICES>
						<CHOICE>AMER</CHOICE>
						<CHOICE>EMEA</CHOICE>
						<CHOICE>APAC</CHOICE>
					</CHOICES>
				</Field>`,
			},
			{
				Name:      "Revenue",
				SchemaXML: `<Field Type="Currency" DisplayName="Revenue, USD" Required="FALSE" LCID="1033" />`,
			},
			{
				Name:      "Manager",
				SchemaXML: `<Field Type="Lookup" DisplayName="Primary Manager" Required="FALSE" List="Lists/Managers" ShowField="Title" />`,
			},
		},
	},
	{
		Title: "Orders",
		Fields: []FieldModel{
			{
				Name:      "Customer",
				SchemaXML: `<Field Type="Lookup" DisplayName="Customer" Required="FALSE" List="Lists/Customers" ShowField="Title" />`,
			},
			{
				Name:      "OrderNumber",
				SchemaXML: `<Field Type="Text" DisplayName="Order Number" Required="FALSE" MaxLength="12" />`,
			},
			{
				Name:      "OrderDate",
				SchemaXML: `<Field Type="DateTime" DisplayName="Order Date" Required="FALSE" />`,
			},
			{
				Name:      "Total",
				SchemaXML: `<Field Type="Currency" DisplayName="Deal Amount, USD" Required="FALSE" LCID="1033" />`,
			},
		},
	},
}
