output "main-vpc-id" {
  value = "${aws_vpc.main.id}"
  description = "The main VPC id"
}

output "public-subnet-id" {
  value = "${aws_subnet.public.id}"
  description = "The public subnet id"
}

output "private-subnet-id" {
  value = "${aws_subnet.private.id}"
  description = "The private subnet id"
}
