package main
import "fmt"

type PaymentProcessor interface{
	Pay(amount float64) error
}

type UPIPayment struct{}

func (u UPIPayment) Pay(amount float64) error{
	fmt.Println("Upi payment done of ",amount)
	return nil
}


type CreditCardPayment struct{}

func (c CreditCardPayment) Pay(amount float64) error{
	fmt.Println("CreditCard payment done of ",amount)
	return nil
}

func ProcessPayment(
    processor PaymentProcessor,
    amount float64,
) error {
    return processor.Pay(amount)
}

func main() {
    upi := UPIPayment{}

    err := ProcessPayment(upi, 500)

    if err != nil {
        fmt.Println(err)
    }

    // You can also pass CreditCardPayment to the same function!
    cc := CreditCardPayment{}
    _ = ProcessPayment(cc, 1500)
}



/*
====================================================================================================
Summary: Java vs. Go Interfaces
====================================================================================================
| Feature                       | Java                                                 | Go                                                            |
|-------------------------------|------------------------------------------------------|---------------------------------------------------------------|
| Declaration                   | interface PaymentProcessor { ... }                   | type PaymentProcessor interface { ... }                       |
| Implementation                | Explicit: class UPI implements PaymentProcessor      | Implicit: Just write the method func (u UPI) Pay(...)         |
| Coupling                      | Tight: The class must know about the interface.      | Loose: The struct does not need to know the interface exists.  |
| Creating Mock/Test interfaces | Harder if 3rd party class didn't declare implements. | Easy: You can define an interface in your package to wrap it.  |
====================================================================================================
*/