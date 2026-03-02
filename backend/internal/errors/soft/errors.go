package softerrors

// import "encoding/json"

// - Soft Errors - // TODO: servono ancora?

// type SoftErr string

// func (err SoftErr) Error() string {
// 	return string(err)
// }

// // Make SoftError serializable
// func (err SoftErr) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(string(err))
// }

// var (
// // Auth
// SoftErrEmailNotSent = SoftErr("s_email_not_sent") // Error sending email
// )

// // CAUTION: does not work with wrapping (fmt.Errorf("...%w...", err))
// func IsSoftErr(err error) bool {
// 	_, ok := err.(SoftErr)
// 	return ok
// }
