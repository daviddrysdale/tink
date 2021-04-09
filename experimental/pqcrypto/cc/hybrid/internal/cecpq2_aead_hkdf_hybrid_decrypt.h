// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
///////////////////////////////////////////////////////////////////////////////

#ifndef THIRD_PARTY_TINK_EXPERIMENTAL_PQCRYPTO_CC_HYBRID_INTERNAL_CECPQ2_AEAD_HKDF_HYBRID_DECRYPT_H_
#define THIRD_PARTY_TINK_EXPERIMENTAL_PQCRYPTO_CC_HYBRID_INTERNAL_CECPQ2_AEAD_HKDF_HYBRID_DECRYPT_H_

#include <memory>
#include <utility>

#include "tink/hybrid_decrypt.h"
#include "tink/util/statusor.h"
#include "pqcrypto/cc/hybrid/cecpq2_aead_hkdf_dem_helper.h"
#include "pqcrypto/cc/subtle/cecpq2_hkdf_recipient_kem_boringssl.h"
#include "pqcrypto/cc/subtle/cecpq2_subtle_boringssl_util.h"

namespace crypto {
namespace tink {

// CECPQ2 decryption with HKDF-KEM (key encapsulation mechanism) and
// AEAD-DEM (data encapsulation mechanism)
class Cecpq2AeadHkdfHybridDecrypt : public HybridDecrypt {
 public:
  // Returns an HybridDecrypt-primitive that uses the key material
  // given in 'recipient_key'
  static crypto::tink::util::StatusOr<std::unique_ptr<HybridDecrypt>> New(
      const Cecpq2AeadHkdfPrivateKeyInternal& private_key_internal);

  crypto::tink::util::StatusOr<std::string> Decrypt(
      absl::string_view ciphertext,
      absl::string_view context_info) const override;

 private:
  Cecpq2AeadHkdfHybridDecrypt(
      const Cecpq2AeadHkdfParamsInternal& recipient_key_params,
      std::unique_ptr<const subtle::Cecpq2HkdfRecipientKemBoringSsl> kem,
      std::unique_ptr<const Cecpq2AeadHkdfDemHelper> dem_helper)
      : recipient_key_params_(recipient_key_params),
        recipient_kem_(std::move(kem)),
        dem_helper_(std::move(dem_helper)) {}

  Cecpq2AeadHkdfParamsInternal recipient_key_params_;
  std::unique_ptr<const subtle::Cecpq2HkdfRecipientKemBoringSsl> recipient_kem_;
  std::unique_ptr<const Cecpq2AeadHkdfDemHelper> dem_helper_;
};

}  // namespace tink
}  // namespace crypto

#endif  // THIRD_PARTY_TINK_EXPERIMENTAL_PQCRYPTO_CC_HYBRID_INTERNAL_CECPQ2_AEAD_HKDF_HYBRID_DECRYPT_H_