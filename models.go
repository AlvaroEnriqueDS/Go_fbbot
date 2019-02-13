package main

type RequestMessage struct {
        Object string  `json:"object"`
        Entry  []Entry `json:"entry"`
}

type Entry struct {
        ID string `json:"id"`
        Time int64 `json:"time"`
        Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
        Sender `json:"sender"`
        Recipient `json:"recipient"`
        Message `json:"message"`
        Timestamp int64 `json:"timestamp"`
}

type Sender struct {
        ID string `json:"id"`
}
type Recipient struct {
        ID string `json:"id"`
}
type Message struct {
       ID string `json:"mid"`
       Seq int64 `json:"seq"`
       Text string `json:"text"`
}


type ResponseMessage struct {
        Recipient `json:"recipient"`
        MessageContent `json:"message"`
}

type MessageContent struct {
        Text string `json:"text,omitempty"`
}
