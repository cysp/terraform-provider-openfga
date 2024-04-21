terraform {
  required_providers {
    openfga = {}
  }
}

output "model_json" {
  value = provider::openfga::transform_model_to_json(file("model.fga"))
}
