package validate

import (
	"log"
	"testing"
)

func TestValidate(t *testing.T) {
	type User struct {
		FirstName string `validate:"required,min=2,max=64,nonNumericStart"`
		LastName  string `validate:"max=64,nonNumericStart"`
		Username  string `validate:"min=3,nonNumericStart"`
		Email     string `validate:"email,optional"`
		Domain    string `validate:"domain,optional"`
		Age       int    `validate:"optional,min=18"`
	}

	type args struct {
		val interface{}
		cfg *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"testGood",
			args{
				User{
					FirstName: "Robert",
					LastName:  "Tucker",
					Username:  "apple",
					Email:     "r.tucker@chcfl.org",
					Domain:    "chcfl.org",
				},
				&Config{},
			},
			false,
		},
		{
			"test2",
			args{
				User{
					FirstName: "Robert",
					LastName:  "4Tucker",
					Username:  "apple",
					Email:     "roborobs_computers@live.com",
				},
				&Config{},
			},
			true,
		},
		{
			"test3",
			args{
				User{
					FirstName: "A",
					Username:  "apple",
					Email:     "roborobs_computers@live.com",
				},
				&Config{},
			},
			true,
		},
		{
			"test4",
			args{
				User{
					FirstName: "Robert",
					LastName:  "Tucker",
					Username:  "apple",
					Email:     "r.tucker@chcfl.com",
				},
				&Config{},
			},
			false,
		},
		{
			"age Failure",
			args{
				User{
					FirstName: "Robert",
					LastName:  "Tucker",
					Username:  "apple",
					Email:     "r.tucker@chcfl.com",
					Age:       12,
				},
				&Config{},
			},
			true,
		},
		{
			"age Pass",
			args{
				User{
					FirstName: "Robert",
					LastName:  "Tucker",
					Username:  "apple",
					Email:     "r.tucker@chcfl.com",
					Age:       19,
				},
				&Config{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.val, tt.args.cfg); (err != nil) != tt.wantErr {
				log.Println("working on test:", tt.name)
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
