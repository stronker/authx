/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package authx

import "time"

type Config struct {
	Port       int
	Secret     string
	ExpirationTime time.Duration
}


