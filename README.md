# Julio Verne
Teste Julio Verne - Oportunidade Mercado Pago

Meli Júlio Verne Test
O projeto tinha o objetivo de desenvolver uma API com a finalidade de traduzir a língua nativa dos IT's (Intra terrenos) através de um JSON que continha as combinações aleatórias de vogais e consoantes, os valores das letras IT eram representados por ASCII do alfabeto português.

Todo o desenvolvimento foi feito em Golang, por ser uma linguagem mais simples e muito boa para performances de rede, e já tinha os package prontos como o MUX, que executa os roteadores http, e o do MySQL que seria o drive. Também fiz um container no Docker.

Quanto ao banco, eu linkei a API com o Banco de Dados MySQL.

### Hospedagem

O projeto foi hospedado na Amazon AWS e pode ser acessado pelo link: 
[Julio Verne](http://ec2-18-217-160-250.us-east-2.compute.amazonaws.com:3000/)

### Execução

É possível testar a aplicação utilizando o Postman. Os endpoints são:
- http://ec2-18-217-160-250.us-east-2.compute.amazonaws.com:3000/translate - método POST com envio de JSON
- http://ec2-18-217-160-250.us-east-2.compute.amazonaws.com:3000/translate/:word - método GET


### Autor

* **Vinicius Coelho Porto**