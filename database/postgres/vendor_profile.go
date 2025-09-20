package postgres

import (
	"context"
	"qvickly/models/vendors"
)

func GetVendorProfile(phoneNumber, password string) (vendorDetails *vendors.VendorProfileDetails, err error) {
	vendorDetails = new(vendors.VendorProfileDetails)
	err = pgPool.QueryRow(context.Background(), `SELECT vendor_id, image_url, business_name, owner_name, is_active FROM quickkart.profile.vendors where phone=$1 and password=$2`, phoneNumber, password).Scan(&vendorDetails.VendorId, &vendorDetails.ImageS3URL, &vendorDetails.StoreName, &vendorDetails.OwnerName, &vendorDetails.StoreLiveStatus)
	if err != nil {
		vendorDetails = nil
	}

	return
}

func AddVendorProfile(data vendors.CompleteVendorProfile) (err error) {
	_, err = pgPool.Exec(context.Background(), `INSERT INTO quickkart.profile.vendors(phone, password, aadhar, account_type, business_name, owner_name, address, latitude, longitude, gstin_number, opening_time, closing_time, image_url, is_active) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, data.Phone, data.Password, data.Aadhar, data.AccountType, data.BusinessName, data.OwnerName, data.Address, data.Latitude, data.Longitude, data.GSTIN, data.OpeningTime, data.ClosingTime, data.ImageS3URL, data.StoreLiveStatus)

	return
}

func GetProfileVendorStatus(vendorId string) (status bool, err error) {
	err = pgPool.QueryRow(context.Background(), `SELECT is_active FROM quickkart.profile.vendors where vendor_id=$1`, vendorId).Scan(&status)
	if err != nil {
		status = false
	}

	return
}

func SetProfileVendorStatus(vendorId string, status bool) (err error) {
	_, err = pgPool.Exec(context.Background(), `UPDATE quickkart.profile.vendors  set is_active = $2 where vendor_id=$1`, vendorId, status)
	if err != nil {
		return err
	}

	return
}
