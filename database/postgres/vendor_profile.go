package postgres

import (
	"context"
	"qvickly/models/vendors"
)

func GetVendorProfile(phoneNumber, password string) (vendorDetails *vendors.VendorProfileDetails, err error) {
	vendorDetails = new(vendors.VendorProfileDetails)
	err = pgPool.QueryRow(context.Background(), `SELECT id, image_url, business_name, owner_name, live_status FROM postgres.vendor_accounts.vendor_accounts where phone_number=$1 and password=$2`, phoneNumber, password).Scan(&vendorDetails.VendorId, &vendorDetails.ImageS3URL, &vendorDetails.StoreName, &vendorDetails.OwnerName, &vendorDetails.StoreLiveStatus)
	if err != nil {
		vendorDetails = nil
	}

	return
}

func AddVendorProfile(data vendors.CompleteVendorProfile) (err error) {
	_, err = pgPool.Exec(context.Background(), `INSERT INTO postgres.vendor_accounts.vendor_accounts(phone_number, account_type, business_name, owner_name, email, address, latitude, longitude, gstin_number, opening_time, closing_time, image_url, live_status) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`, data.Phone, data.AccountType, data.BusinessName, data.OwnerName, data.Email, data.Address, data.Latitude, data.Longitude, data.GSTIN, data.OpeningTime, data.ClosingTime, data.ImageS3URL, data.StoreLiveStatus)

	return
}

func GetProfileVendorStatus(vendorId string) (status bool, err error) {
	err = pgPool.QueryRow(context.Background(), `SELECT live_status FROM postgres.vendor_accounts.vendor_accounts where id=$1::uuid`, vendorId).Scan(&status)
	if err != nil {
		status = false
	}

	return
}

func SetProfileVendorStatus(vendorId string, status bool) (err error) {
	_, err = pgPool.Exec(context.Background(), `UPDATE postgres.vendor_accounts.vendor_accounts  set live_status = $2 where id=$1::uuid`, vendorId, status)
	if err != nil {
		return err
	}

	return
}
