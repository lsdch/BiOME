import js from '@eslint/js'
import pluginCypress from 'eslint-plugin-cypress'
import pluginVue from 'eslint-plugin-vue'
import {
  defineConfigWithVueTs,
  vueTsConfigs,
} from '@vue/eslint-config-typescript'
import eslintConfigPrettier from '@vue/eslint-config-prettier'

export default defineConfigWithVueTs(
  js.configs.recommended,
  pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,
  eslintConfigPrettier,

  {
    files: ['**/*.{js,jsx,ts,tsx,vue}'],

    languageOptions: {
      ecmaVersion: 'latest',
    },

    rules: {
      'no-unused-vars': 'off',

      'vue/no-unused-vars': [
        'warn',
        {
            ignorePattern: '^_',
        },
    ],

      '@typescript-eslint/no-unused-vars': [
        'warn',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
        },
      ],
    },
  },

  {
    files: ['cypress/e2e/**/*.{cy,spec}.{js,ts,jsx,tsx}'],
    ...pluginCypress.configs.recommended,
  },
)