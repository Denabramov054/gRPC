package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc/pkg/api"
	"log"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewUserClient(conn)

	//Create user method
	log.Println("Create user request")
	var name string
	var email string
	var password string
	log.Print("Enter name: ")
	fmt.Scanln(&name)
	log.Print("Enter email: ")
	fmt.Scanln(&email)
	log.Print("Enter password: ")
	fmt.Scanln(&password)
	res, err := client.Create(context.Background(), &api.CreateRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Create user response: ", res.String())

	//Getting all users method
	log.Println("Getting all users request...")
	usersRes, err := client.Users(context.Background(), &api.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Get users response: ", usersRes.String())

	//Delete user method
	log.Println("Deleting user request...")
	var id uint64
	log.Print("Enter id: ")
	fmt.Scanln(&id)
	delRes, err := client.Delete(context.Background(), &api.DeleteRequest{Id: id})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Delete user response: ", delRes.String())

}
