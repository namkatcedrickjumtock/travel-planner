package travelplanner

import "context"


func CreateOrder(ctx context.Context, req OrderReq) error {
    user, err := users.Get(ctx, req.UserID)
    if err != nil {
        return fmt.Errorf("couldn’t get user by id. why=%w", err)
    }

    err = inventory.Reserve(ctx, req.ItemID, req.Qty)
    if err != nil {
        return fmt.Errorf("error reserving item=%s. why=%w", req.ItemID, err)
    }

    err = payments.Charge(ctx, user.PaymentID, req.Total)
    if err != nil {
        return fmt.Errorf("error charging item=%s. why=%w", req.ItemID, err)
    }

    return saveOrder(ctx, req.UserID, req.ItemID)
}