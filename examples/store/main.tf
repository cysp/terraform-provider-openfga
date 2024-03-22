terraform {
  required_providers {
    openfga = {
      source = "cysp/openfga"
    }
  }
}

provider "openfga" {
  api_url = "http://localhost:8080"
}

resource "openfga_store" "this" {
  name = "example"
}

# resource "openfga_store" "planner_iq" {
#   name = "planner-iq"
# }

# resource "openfga_authorization_model" "planner_iq" {
#   store_id = openfga_store.planner_iq.id

#   model_json = file("files/authorization_model.json")
# }

# import {
#   id = "01HQ85VP9TPSEKA49MFJ0R9D9H"
#   to = openfga_store.planner_iq
# }

# import {
#   id = "01HQ85VP9TPSEKA49MFJ0R9D9H/01HQ85XXB7B3110333Z3ZDG7ZJ"
#   to = openfga_authorization_model.planner_iq
# }

# resource "openfga_tuple" "planner_iq_foo" {
#   store_id = openfga_store.planner_iq.id

#   user     = "role:editor#assignee"
#   relation = "editor"
#   object   = "campaign:12"
# }
