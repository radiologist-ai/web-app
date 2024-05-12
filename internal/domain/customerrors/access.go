package customerrors

import "fmt"

var (
	AccessError                = fmt.Errorf("access error. ")
	AccessErrorNotEnoughRights = fmt.Errorf("%wnot enough rights for operation. ", AccessError)
	NeedToBeDoctor             = fmt.Errorf("%wonly doctor can perform this operation. ", AccessErrorNotEnoughRights)
	NeedToBePatient            = fmt.Errorf("%wonly patient can perform this operation. ", AccessErrorNotEnoughRights)
)
