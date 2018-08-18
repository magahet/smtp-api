package main

import "testing"

func TestSender_Send(t *testing.T) {
	type fields struct {
		host string
		port int
	}
	type args struct {
		message *Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "simple",
			fields:  fields{"localhost", 1025},
			args:    args{&Message{"from@blah.com", "testing", "blah blah", []string{"to@blah.com"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sender{
				host: tt.fields.host,
				port: tt.fields.port,
			}
			if err := s.Send(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Sender.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
