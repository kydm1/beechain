package executor

import "fmt"

var (
	goldbillPre = "mavl-goldbill-"
	localGoldbillPre = "LODB-goldbill-"
)

func calcGoldbillPlatformKey() []byte {
	return []byte(fmt.Sprintf(goldbillPre+"platform-"+"%s","platform"))
}

func calcGoldbillUserKey(addr string) []byte {
	return []byte(fmt.Sprintf(goldbillPre+"user-"+"%s",addr))
}


func calcGoldbillAdminKey(addr string) []byte {
	return []byte(fmt.Sprintf(goldbillPre+"%s",addr))
}

func calcGoldbillDetailKey(id string) []byte {
	return []byte(fmt.Sprintf(goldbillPre+"bill-"+"%s",id))
}

func calcGoldbillInvoice(id string) []byte {
	return []byte(fmt.Sprintf(goldbillPre+"invoice"+"%s",id))
}

func calcGoldbillVATInvoice(id string) []byte {
	return []byte(fmt.Sprintf(goldbillPre+"vatinvoice"+"%s",id))
}

func calcGoldbillUserState() []byte {
	return []byte(localGoldbillPre+"userstate")
}
