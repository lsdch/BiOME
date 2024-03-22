<template>
  <v-form>
    <v-row>
      <v-col>
        <v-card
          class="flex-grow-1"
          :max-width="600"
          subtitle="Minimum password strength"
          variant="text"
        >
          <v-card-text>
            <v-slider
              :step="1"
              :min="3"
              :max="5"
              :ticks="{ 3: 'Medium', 4: 'Strong', 5: 'Very strong' }"
              show-ticks="always"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6">
        <v-text-field
          class="flex-grow-1"
          label="Authentication token secret key"
          :min="32"
          prepend-inner-icon="mdi-key"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <VNumberInput
          class="flex-grow-0"
          label="User session lifetime (seconds)"
          :step="60"
          :min="60"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <VNumberInput label="Account token lifetime (hours)" :step="1" :min="1" />
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts"></script>

<style scoped></style>
