Feature: Admin Configuration panel
	create admin init config 
	get config and try to create another config

  Scenario: Create admin config
     Given admin config data 
     When is created
     AND get config
     Then try to create a new config
   
