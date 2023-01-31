1. Make build the provider
2. https://www.infracloud.io/blogs/developing-terraform-custom-provider/ Rename the provider and copy it
3. Run a dockerfile with sonarqube as an image, change password to admin1 before running tests
4. Run the code in examples/basic

# TODO

* Improve the manual check process, also the importing of the goddamn provider is mad ugly
* Revert all our messy changes to files just to test
* Perhaps make it so that the addon checks the default branch of the system so we can revert it back to the actual
  default