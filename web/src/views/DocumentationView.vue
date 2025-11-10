<template>
  <div class="docs">
    <nav class="docs-nav">
      <h3>Documentation</h3>
      <ul>
        <li v-for="doc in docs" :key="doc.id">
          <a
            :href="`#${doc.id}`"
            :class="{ active: activeSection === doc.id }"
            @click.prevent="scrollToSection(doc.id)"
          >
            {{ doc.title }}
          </a>
        </li>
      </ul>
    </nav>

    <main class="docs-content">
      <div v-if="docs.length === 0" class="loading">
        Loading documentation...
      </div>
      <section
        v-for="doc in docs"
        v-else
        :key="doc.id"
        :id="doc.id"
        class="docs-section"
      >
        <h2>{{ doc.title }}</h2>
        <div class="markdown-body" v-html="doc.html"></div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { marked } from "marked";

interface Doc {
  id: string;
  title: string;
  html: string;
}

const docs = ref<Doc[]>([]);
const activeSection = ref<string>("");

onMounted(async () => {
  // Load markdown files from public/docs folder
  const docFiles = ["getting-started", "guide", "api"];

  for (const file of docFiles) {
    try {
      const url = `/docs/${file}.md`;
      console.log(`Fetching: ${url}`);
      const response = await fetch(url);
      console.log(`Response for ${file}.md:`, response.status);
      if (response.ok) {
        const markdown = await response.text();
        const html = await marked(markdown);
        docs.value.push({
          id: file,
          title: formatTitle(file),
          html,
        });
      } else {
        console.warn(`Failed to fetch ${file}.md: ${response.status}`);
      }
    } catch (error) {
      console.error(`Error loading ${file}.md:`, error);
    }
  }
});

function formatTitle(filename: string): string {
  return filename
    .split("-")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

function scrollToSection(id: string) {
  activeSection.value = id;
  const element = document.getElementById(id);
  if (element) {
    element.scrollIntoView({ behavior: "smooth" });
  }
}
</script>

<style scoped>
.docs {
  display: grid;
  grid-template-columns: 250px 1fr;
  gap: var(--spacing-xl);
  max-width: 1200px;
  margin: 0 auto;
  align-items: start;
}

.docs-nav {
  position: sticky;
  top: var(--spacing-sm);
  height: fit-content;
  padding: var(--spacing-md);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  background: var(--bg);
}

.docs-nav h3 {
  margin: 0 0 var(--spacing-md) 0;
  font-size: 0.9375rem;
  color: var(--fg-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.docs-nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.docs-nav a {
  display: block;
  padding: 0.375rem 0.75rem;
  border-radius: var(--radius-sm);
  color: var(--fg-secondary);
  text-decoration: none;
  transition: var(--transition);
  font-size: 0.875rem;
}

.docs-nav a:hover {
  color: var(--accent);
  background: transparent;
}

.docs-nav a.active {
  color: var(--accent);
  background: transparent;
  font-weight: 500;
}

.docs-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xl);
}

.docs-section {
  scroll-margin-top: 100px;
  padding-bottom: var(--spacing-xl);
  border-bottom: 1px solid var(--border);
}

.docs-section:last-child {
  border-bottom: none;
}

.docs-section h2 {
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border);
  margin-bottom: var(--spacing-xl);
}

.markdown-body {
  line-height: 1.7;
  color: var(--fg);
  user-select: text;
  -webkit-user-select: text;
}

.markdown-body h3 {
  margin-top: var(--spacing-lg);
  margin-bottom: var(--spacing-md);
  font-size: 1.125rem;
}

.markdown-body h4 {
  margin-top: var(--spacing-md);
  margin-bottom: var(--spacing-sm);
  font-size: 1rem;
}

.markdown-body p {
  margin: var(--spacing-md) 0;
}

.markdown-body ul,
.markdown-body ol {
  margin: var(--spacing-md) 0;
  padding-left: 1.5rem;
}

.markdown-body li {
  margin: var(--spacing-sm) 0;
}

.markdown-body code {
  background: var(--bg);
  color: var(--accent);
  padding: 0.125rem 0.375rem;
  border-radius: var(--radius-sm);
  font-family: "JetBrains Mono", "Roboto Mono", monospace;
  font-size: 0.875em;
}

.markdown-body pre {
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  overflow-x: auto;
  margin: var(--spacing-md) 0;
}

.markdown-body pre code {
  background: transparent;
  color: var(--fg);
  padding: 0;
}

.markdown-body blockquote {
  border-left: 3px solid var(--accent);
  padding-left: var(--spacing-md);
  margin: var(--spacing-md) 0;
  color: var(--fg-secondary);
  font-style: italic;
}

.markdown-body a {
  color: var(--accent);
  text-decoration: none;
  border-bottom: 1px dotted var(--accent);
  transition: var(--transition);
}

.markdown-body a:hover {
  opacity: 0.8;
}

.markdown-body table {
  width: 100%;
  border-collapse: collapse;
  margin: var(--spacing-md) 0;
}

.markdown-body th,
.markdown-body td {
  padding: var(--spacing-sm);
  border: 1px solid var(--border);
  text-align: left;
}

.markdown-body th {
  background: transparent;
  font-weight: 600;
  color: var(--accent);
}

@media (max-width: 768px) {
  .docs {
    grid-template-columns: 1fr;
  }

  .docs-nav {
    position: static;
    max-height: none;
    border: none;
    padding: 0;
    background: transparent;
  }

  .docs-nav h3 {
    display: none;
  }

  .docs-nav ul {
    flex-direction: row;
    flex-wrap: wrap;
    gap: var(--spacing-md);
  }

  .docs-nav a {
    padding: 0.375rem 0.75rem;
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
  }
}

.loading {
  text-align: center;
  color: var(--fg-secondary);
  padding: var(--spacing-lg);
  font-style: italic;
}
</style>
