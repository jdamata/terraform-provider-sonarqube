variable "name" {
  type        = string
  description = "Name of project"
}

variable "quality_gates" {
  type        = map(any)
  description = "Quality gates"
}