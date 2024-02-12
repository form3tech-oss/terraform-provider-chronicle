provider "chronicle" {
  backstoryapi_credentials = "backstoryapi_crendentials_json_string"
  region                   = "europe"
  request_attempts         = 5
  request_timeout          = 120
}
