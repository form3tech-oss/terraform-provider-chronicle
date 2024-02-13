resource "chronicle_reference_list" "list" {
  name         = "mylist"
  description  = "my awesome list"
  content_type = "CONTENT_TYPE_DEFAULT_STRING"
  lines = [
    "one",
    "two"
  ]
}
