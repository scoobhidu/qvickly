package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var pgPool *pgxpool.Pool

func InitConn() (err error) {
	connConfig := "user=" + os.Getenv("user") +
		" password=" + os.Getenv("password") +
		" host=" + os.Getenv("host") +
		" port=" + os.Getenv("port") +
		" dbname=" + os.Getenv("dbname")
	// postgres://jack:secret@foo.example.com:5432,bar.example.com:5432/mydb
	// context.Background gives an empty context, don't want to give it a time to close automatically
	// if the connection does not happen within t time
	pgPool, err = pgxpool.New(context.Background(), connConfig)

	if err != nil {
		log.Panicf("Unable to connect to database: %v\n", err) // printf followed by panic(), no need for os.Exit as panic will exit
	}
	return err
}

func CloseConn() {
	defer pgPool.Close()
}

//const (
//	// paths of what client is placing on forward path
//	DELIVERY string = "/orders/v1/delivery"
//	INTRADAY string = "/orders/v1/intraday/regular"
//
//	BRACKET string = "/orders/v1/intraday/bracket"
//	COVER   string = "/orders/v1/intraday/cover"
//
//	BASKET string = "/orders/v1/basket"
//
//	CXL    string = "/orders/v1/cancel"
//	MODIFY string = "/orders/v1/modify"
//
//	// actions made by the client on their order
//	ACTION_TYPE_NEW    string = "NEW"
//	ACTION_TYPE_CXL    string = "CXL"
//	ACTION_TYPE_MODIFY string = "MODIFY"
//)

//func WriteData(event models.LambdaEvent) (err error) {
//	if event.Path == "" {
//		err = errors.New("the event does not contain path")
//		return
//	}
//
//	if event.Path == DELIVERY || event.Path == INTRADAY || event.Path == COVER || event.Path == BRACKET {
//		err = newOrder(event)
//	} else if event.Path == BASKET {
//		err = basketOrder(event)
//	} else if event.Path == CXL {
//		err = cxlOrder(event)
//	} else if event.Path == MODIFY {
//		err = modifyOrder(event)
//	}
//
//	return
//}

// new order event insertion
//func newOrder(event models.LambdaEvent) (err error) {
//	// off_market_flag by default '0' because the table column expects a char
//	// if we get the off market flag from SNS event and it is 1 only then offMktFlag will be set to '1'
//	offMktFlag := "'0'"
//	if event.OffMktFlag != nil && *(event.OffMktFlag) == 1 {
//		offMktFlag = "'1'"
//	}
//
//	query := ""
//
//	// if we have received status in the event then it's a response from RS and not the initial order request from fwd path
//	if event.Status != nil {
//		// the event.OmsResponses is in the struct format which gets converted to a bytes format (omsResponse)
//		// and those bytes are then passed as string to insert in jsonb column of oms_response
//		omsResponse, err := json.Marshal(event.OmsResponses)
//		if err != nil {
//			return errors.New("Unable to marshal omsResponse to string format" + err.Error())
//		}
//
//		query = fmt.Sprintf(
//			`insert into txns.order_details(security_id, request_uuid, client_id, qty, disc_qty, status, remarks,
//			    t1, t2, t3, t4, action_type, oms_response, rs_orderid, epoch, message)
//		   		values(%v, '%v', '%v', %v, %v, %v, %v, %v, %v, %v, %v, '%v', '%v', %v, %v, %v);`,
//
//			utils.GetNullableInt64(event.SecurityId),
//			event.Uuid,
//			event.ClientID,
//			utils.GetNullableInt64(event.Qty),
//			utils.GetNullableInt64(event.DiscQty),
//			utils.GetNullableString(event.Status),
//			utils.GetNullableString(event.Remarks),
//			utils.GetNullableInt64(event.T1),
//			utils.GetNullableInt64(event.T2),
//			utils.GetNullableInt64(event.T3),
//			utils.GetNullableInt64(event.T4),
//			ACTION_TYPE_NEW,
//			string(omsResponse), // byte format is converted to string
//			utils.ConvertNullableStringToInt64(event.OmsOrderId),
//			time.Now().Unix(),
//			utils.GetNullableString(event.Message),
//		)
//	} else {
//		// the initial request from fwd path
//		query = fmt.Sprintf(
//			`insert into txns.order_details(request_uuid, client_id, security_id,
//			   exchange, segment, txn_type, order_type, qty, price, trigger_price, product, validity, disc_qty,
//			   off_mkt_flag, remarks, target_price, stoploss_price, ltp, epoch, action_type)
//			   values('%v', '%v', %v, '%v', '%v', '%v', %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, '%v');`,
//
//			event.Uuid,
//			event.ClientID,
//			utils.GetNullableInt64(event.SecurityId),
//			event.Exchange,
//			event.Segment,
//			event.TxnType,
//			utils.GetNullableString(event.OrderType),
//			utils.GetNullableInt64(event.Qty),
//			utils.GetNullableFloat32(event.Price),
//			utils.GetNullableFloat32(event.TriggerPrice),
//			utils.GetNullableString(event.Product),
//			utils.GetNullableString(event.Validity),
//			utils.GetNullableInt64(event.DiscQty),
//			offMktFlag,
//			utils.GetNullableString(event.Remarks),
//			utils.GetNullableFloat32(event.TargetPrice),
//			utils.GetNullableFloat32(event.StoplossPrice),
//			utils.GetNullableFloat32(event.Ltp),
//			time.Now().Unix(),
//			ACTION_TYPE_NEW,
//		)
//	}
//
//	log.Println(query)
//
//	_, err = pgPool.Exec(context.Background(), query)
//	if err != nil {
//		log.Println("There was an error executing this query" + err.Error())
//		return
//	}
//
//	return
//}

// new basket order event insertion
//func basketOrder(event models.LambdaEvent) (err error) {
//	// if we have received status in the event then it's a response from RS and not the initial order request from fwd path
//	if event.Status != nil {
//
//		// if no basket child was received from the RS response then there is nothing to insert to table
//		if event.BasketData == nil {
//			log.Println("empty basket response found")
//		} else {
//			for _, basketResponse := range *event.BasketData {
//
//				// the basketResponse.Response.OmsRes (1st element of every basket child response) - is in the struct format
//				// which gets converted to a bytes format (omsResponse)
//				// and those bytes are then passed as string to insert in jsonb column of oms_response
//				omsResponse, err := json.Marshal(basketResponse.Response.OmsRes)
//				if err != nil {
//					return errors.New("Unable to marshal omsResponse to string format" + err.Error())
//				}
//
//				query := fmt.Sprintf(
//					`insert into txns.order_details(security_id, request_uuid, client_id, qty, disc_qty, status, remarks,
//			    t1, t2, t3, t4, action_type, oms_response, rs_orderid, epoch, message)
//		   		values(%v, '%v', '%v', %v, %v, %v, %v, %v, %v, %v, %v, '%v', '%v', %v, %v, %v);`,
//
//					utils.GetNullableInt64(basketResponse.Response.SecurityId),
//					event.Uuid,
//					event.ClientID,
//					utils.GetNullableInt64(basketResponse.Response.Qty),
//					utils.GetNullableInt64(basketResponse.Response.DiscQty),
//					utils.GetNullableString(basketResponse.Response.Status),
//					utils.GetNullableString(basketResponse.Response.Remarks),
//					utils.GetNullableInt64(basketResponse.Response.T1),
//					utils.GetNullableInt64(basketResponse.Response.T2),
//					utils.GetNullableInt64(basketResponse.Response.T3),
//					utils.GetNullableInt64(basketResponse.Response.T4),
//					ACTION_TYPE_NEW,
//					string(omsResponse), // byte format is converted to string
//					utils.ConvertNullableStringToInt64(basketResponse.Response.OmsOrderId),
//					time.Now().Unix(),
//					utils.GetNullableString(basketResponse.Response.Message),
//				)
//
//				log.Println(query)
//
//				_, err = pgPool.Exec(context.Background(), query)
//				if err != nil {
//					log.Println("There was an error executing this query" + err.Error())
//				}
//			}
//		}
//	} else {
//		// the initial request from fwd path
//		query := fmt.Sprintf(
//			`insert into txns.order_details(request_uuid, client_id, epoch, action_type)
//			   values('%v', '%v', %v, '%v');`,
//
//			event.Uuid,
//			event.ClientID,
//			time.Now().Unix(),
//			ACTION_TYPE_NEW,
//		)
//
//		log.Println(query)
//
//		_, err = pgPool.Exec(context.Background(), query)
//		if err != nil {
//			log.Println("There was an error executing this query" + err.Error())
//			return
//		}
//	}
//	return
//}

// cancel order event insertion
//func cxlOrder(event models.LambdaEvent) (err error) {
//	query := ""
//
//	// if we have received status in the event then it's a response from RS and not the initial order request from fwd path
//	if event.Status != nil {
//		// the event.OmsResponses is in the struct format which gets converted to a bytes format (omsResponse)
//		// and those bytes are then passed as string to insert in jsonb column of oms_response
//		omsResponse, err := json.Marshal(event.OmsResponses)
//		if err != nil {
//			return errors.New("Unable to marshal omsResponse to string format" + err.Error())
//		}
//
//		query = fmt.Sprintf(
//			`insert into txns.order_details(request_uuid, client_id, oms_response, message, status, epoch)
//				values('%v', '%v', '%v', %v, %v, %v)`,
//
//			event.Uuid,
//			event.ClientID,
//			string(omsResponse),
//			utils.GetNullableString(event.Message),
//			utils.GetNullableString(event.Status),
//			time.Now().Unix(),
//		)
//	} else {
//		// the initial order request from fwd path
//		query = fmt.Sprintf(
//			`insert into txns.order_details(request_uuid, client_id, rs_orderid, epoch, action_type)
//		    	values('%v', '%v', %v, %v, '%v')`,
//
//			event.Uuid,
//			event.ClientID,
//			utils.ConvertNullableStringToInt64(event.OrderNo), // the order no is expected in int format
//			time.Now().Unix(),
//			ACTION_TYPE_CXL,
//		)
//	}
//
//	log.Println(query)
//
//	_, err = pgPool.Exec(context.Background(), query)
//	if err != nil {
//		log.Println("There was an error executing this query" + err.Error())
//		return
//	}
//
//	return
//}

// modify order event insertion
//func modifyOrder(event models.LambdaEvent) (err error) {
//	query := ""
//
//	// if we have received status in the event then it's a response from RS and not the initial order request from fwd path
//	if event.Status != nil {
//		// the event.OmsResponses is in the struct format which gets converted to a bytes format (omsResponse)
//		// and those bytes are then passed as string to insert in jsonb column of oms_response
//		omsResponse, err := json.Marshal(event.OmsResponses)
//		if err != nil {
//			return errors.New("Unable to marshal omsResponse to string format" + err.Error())
//		}
//
//		query = fmt.Sprintf(
//			`insert into txns.order_details(request_uuid, client_id, oms_response, message, status, epoch, action_type)
//				values('%v', '%v', '%v', %v, %v, %v, '%v')`,
//
//			event.Uuid,
//			event.ClientID,
//			string(omsResponse),
//			utils.GetNullableString(event.Message),
//			utils.GetNullableString(event.Status),
//			time.Now().Unix(),
//			ACTION_TYPE_MODIFY,
//		)
//	} else {
//		// the initial order request from fwd path
//		offMktFlag := "'0'"
//		if event.OffMktFlag != nil && *(event.OffMktFlag) == 1 {
//			offMktFlag = "'1'"
//		}
//
//		query = fmt.Sprintf(
//			`insert into txns.order_details(request_uuid, client_id, rs_orderid, order_type, qty, price, trigger_price, validity,
//			  disc_qty, off_mkt_flag, remarks, target_price, stoploss_price, epoch, action_type)
//		      values('%v', '%v', %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, '%v')`,
//
//			event.Uuid,
//			event.ClientID,
//			utils.GetNullableString(event.OrderNo),
//			utils.GetNullableString(event.OrderType),
//			utils.GetNullableInt64(event.Qty),
//			utils.GetNullableFloat32(event.Price),
//			utils.GetNullableFloat32(event.TriggerPrice),
//			utils.GetNullableString(event.Validity),
//			utils.GetNullableInt64(event.DiscQty),
//			offMktFlag,
//			utils.GetNullableString(event.Remarks),
//			utils.GetNullableFloat32(event.TargetPrice),
//			utils.GetNullableFloat32(event.StoplossPrice),
//			time.Now().Unix(),
//			ACTION_TYPE_MODIFY,
//		)
//	}
//
//	log.Println(query)
//
//	_, err = pgPool.Exec(context.Background(), query)
//	if err != nil {
//		log.Println("There was an error executing this query" + err.Error())
//		return
//	}
//
//	return
//}
