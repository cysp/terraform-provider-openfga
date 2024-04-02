

output "test" {
  value = provider::openfga::transform_model_to_json(file("model.fga"))
}
