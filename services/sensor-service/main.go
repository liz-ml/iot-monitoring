package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SensorCreated struct {
	Event    string  `json:"event"`
	SensorID string  `json:"sensor_id"`
	Name     string  `json:"name"`
	Limit    float64 `json:"limit"`
}

func main() {
	// 1. Conectar ao RabbitMQ
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL não configurada")
	}

	conn, err := amqp.Dial(url) //estabelece conexão com o servidor RabbitMQ
	//forma simplificada de sua assinatura é: func Dial(url string) (*Connection, error)
	if err != nil {
		//log.Fatal não é uma boa escolha em um handler HTTP, pq encerra o microsserviço inteiro
		log.Fatal("Erro na conexão: ", err)
	}
	defer conn.Close()

	// 2. Abrir um canal AMQP; um channel é um canal lógico de comunicação
	// dentro dessa conexão para realizar operações AMQP
	ch, err := conn.Channel() //método do objeto conn ; diferente da função do pacote amqp -> amqp.Dial(url)
	if err != nil {
		log.Fatal("Erro ao abrir canal: ", err)
	}
	defer ch.Close()

	// 3. Declarar uma exchange
	err = ch.ExchangeDeclare(
		//os argumentos são name string, kind string, durable bool, autoDelete bool, internal bool;
		//noWait bool, args amqp.Table
		"sensor.events",
		"topic",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Erro ao criar exchange: ", err)
	}

	// 4. Declarar uma fila de teste
	_, err = ch.QueueDeclare(
		"sensor.test",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Erro ao criar fila: ", err)
	}

	// 5. Associar a fila à exchange
	err = ch.QueueBind(
		"sensor.test",
		"sensor.created",
		"sensor.events",
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Erro ao associar fila: ", err)
	}

	// 6. Criar um evento fictício
	event := SensorCreated{
		Event:    "sensor.created",
		SensorID: "sensor-001",
		Name:     "Temperatura Fortaleza",
		Limit:    35,
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Fatal("Erro ao converter evento: ", err)
	}

	// 7. Publicar mensagem
	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"sensor.events",
		"sensor.created",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		log.Fatal("Erro ao publicar evento: ", err)
	}

	log.Println("Evento enviado:", string(body))
}
