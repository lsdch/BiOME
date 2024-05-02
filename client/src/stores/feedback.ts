import { defineStore } from "pinia";
import { ref } from "vue";

export type FeedbackType = "primary" | "success" | "warning" | "error" | "info"

export type Feedback = {
  type?: FeedbackType,
  message: string
}

export const useFeedback = defineStore("feedback", () => {
  const feedbackQueue = ref<Feedback[]>([])

  function feedback(f: Feedback) {
    if (feedbackQueue.value.length === 0 || feedbackQueue.value[feedbackQueue.value.length - 1] !== f)
      feedbackQueue.value.push(f)

    if (!current.value) next()
  }

  function next() {
    current.value = feedbackQueue.value.shift()
  }

  const current = ref<Feedback | undefined>(undefined)

  return { feedback, next, current }
})