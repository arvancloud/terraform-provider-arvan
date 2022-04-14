version = 0.1.0
provider_path = arvancloud.com/terraform/arvan/$(version)
binary_name = terraform-provider-arvan

clean:
	rm -rf ~/.terraform.d/plugins/$(provider_path)/linux_amd64

install_linux: clean
	mkdir -p ~/.terraform.d/plugins/$(provider_path)/linux_amd64
	go build -o ~/.terraform.d/plugins/$(provider_path)/linux_amd64/$(binary_name)_$(version)
