package validate

import "testing"

func TestValidate(t *testing.T) {
	type User struct {
		Username string `validate:"min=3,nonNumericStart"`
		Email    string `validate:"email,required"`
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
			"test1",
			args{
				User{
					Username: "4pple",
				},
				&Config{},
			},
			true,
		},
		{
			name: "test2",
			args: args{
				val: User{
					Username: "apple",
					Email:    "roborobs_computers@live.com",
				},
				cfg: &Config{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.val, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
