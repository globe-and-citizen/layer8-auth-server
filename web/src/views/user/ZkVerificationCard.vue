<template>
  <div class="verification-card">
    <div class="verification-content">
      <div class="verification-info">
        <div class="verification-title">
          <h3>{{ label }}</h3>

          <span
            v-if="verified"
            class="verification-status verified"
          >
            <span>✓</span>
            <span>Verified</span>
          </span>

          <span
            v-else
            class="verification-status not-verified"
          >
            Not verified
          </span>

          <span
            v-if="location && verified"
            class="verification-location"
          >
            <span>📍</span>
            <span>{{ location }}</span>
          </span>
        </div>

        <p v-if="verified && lastVerifiedAt" class="last-verified">
          Last verified {{ lastVerifiedAt }}
        </p>
      </div>

      <button
        @click="$emit('verify')"
        class="verify-button"
      >
        {{ verifyLabel }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  verified: boolean
  location?: string | null
  lastVerifiedAt?: string | null
  verifyLabel: string
}>()

defineEmits<{
  verify: []
}>()
</script>

<style scoped>
.verification-card {
  border: 1px solid #e5e7eb;
  border-radius: 0.75rem;
  background: white;
  padding: 1.5rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.verification-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
}

.verification-info {
  flex: 1;
  min-width: 0;
}

.verification-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.verification-title h3 {
  margin: 0;
  font-size: 1.25rem;
  line-height: 1.75rem;
  font-weight: 500;
  color: #111827;
}

.verification-status {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.25rem 0.75rem;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
}

.verification-location {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.75rem;
  border-radius: 0.375rem;
  background: #f3f4f6;
  color: #6b7280;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
}

.verification-status.verified {
  background: #dcfce7;
  color: #16a34a;
}

.verification-status.not-verified {
  background: #f3f4f6;
  color: #6b7280;
}

.last-verified {
  margin: 0.5rem 0 0;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #6b7280;
}

.verify-button {
  width: 20rem;
  flex-shrink: 0;
  padding: 1rem 1.5rem;
  border: 2px solid #4f80e1;
  border-radius: 0.5rem;
  background: white;
  color: #4f80e1;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 200ms ease,
  color 200ms ease;
}

.verify-button:hover {
  background: #4f80e1;
  color: white;
}

@media (max-width: 767px) {
  .verification-content {
    flex-direction: column;
    align-items: stretch;
    gap: 1.25rem;
  }

  .verify-button {
    width: 100%;
  }
}
</style>
