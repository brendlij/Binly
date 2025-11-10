import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import PasteView from "../views/PasteView.vue";
import DocumentationView from "../views/DocumentationView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", name: "home", component: HomeView },
    { path: "/p/:id", name: "paste", component: PasteView },
    { path: "/docs", name: "docs", component: DocumentationView },
  ],
  scrollBehavior: () => ({ top: 0 }),
});

export default router;
