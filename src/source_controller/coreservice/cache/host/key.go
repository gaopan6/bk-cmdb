/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package host

import (
	"math/rand"
	"strconv"
	"time"

	"configcenter/src/common"
)

const hostKeyNamespace = common.BKCacheKeyV3Prefix + "host"

var hostKey = hostKeyGenerator{
	namespace: hostKeyNamespace,
	// 30 minutes
	expireSeconds:      30 * 60 * time.Second,
	expireRangeSeconds: [2]int{-600, 600},
}

type hostKeyGenerator struct {
	namespace string
	// expireSeconds is defined how long is the ttl for the key.
	// it's always used with the expireRangeSeconds to avoid the keys is expired at same time, which
	// will have large numbers of request flood to the mongodb, we can not accept that.
	// for example, if expireSeconds is 30min, expireRangeSeconds is [-600, 600], then
	// a key's expire seconds is between [20, 40] minutes.
	expireSeconds time.Duration
	// min:[0], max:[1]
	expireRangeSeconds [2]int
}

func (h hostKeyGenerator) HostDetailKey(hostID int64) string {
	return h.namespace + ":detail:" + strconv.FormatInt(hostID, 10)
}

func (h hostKeyGenerator) HostDetailLockKey(hostID int64) string {
	return h.namespace + ":detail:lock:" + strconv.FormatInt(hostID, 10)
}

// key to store the relation with ip and host id:
// key: bk_host_innerip:bk_cloud_id
// value: bk_host_id
// this key has a ttl, which is h.expireSeconds
func (h hostKeyGenerator) IPCloudIDKey(ip string, cloudID int64) string {
	return h.namespace + ":ip_cloud_id:" + ip + ":" + strconv.FormatInt(cloudID, 10)
}

func (h hostKeyGenerator) ListDoneKey() string {
	return h.namespace + ":listdone"
}

func (h hostKeyGenerator) WithRandomExpireSeconds() time.Duration {
	rand.Seed(time.Now().UnixNano())
	seconds := rand.Intn(h.expireRangeSeconds[1]-h.expireRangeSeconds[0]) + h.expireRangeSeconds[0]
	return h.expireSeconds + time.Duration(seconds)*time.Second
}
