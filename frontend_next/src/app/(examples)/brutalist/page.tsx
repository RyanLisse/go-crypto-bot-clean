import React from "react";
import { BrutalistButton } from "@/components/examples/brutalist-button";

export default function BrutalistDemoPage() {
  return (
    <div className="container mx-auto py-12 space-y-12">
      <div className="space-y-6">
        <h1 className="text-4xl font-bold">Brutalist Design System</h1>
        <p className="text-xl">
          This page showcases the BrutalistButton component with different variants,
          sizes, and interactive styles.
        </p>
      </div>

      <section className="space-y-4">
        <h2 className="text-2xl font-bold">Button Variants</h2>
        <div className="flex flex-wrap gap-4">
          <BrutalistButton>Default</BrutalistButton>
          <BrutalistButton variant="primary">Primary</BrutalistButton>
          <BrutalistButton variant="secondary">Secondary</BrutalistButton>
          <BrutalistButton variant="success">Success</BrutalistButton>
          <BrutalistButton variant="warning">Warning</BrutalistButton>
          <BrutalistButton variant="danger">Danger</BrutalistButton>
          <BrutalistButton variant="ghost">Ghost</BrutalistButton>
          <BrutalistButton variant="link">Link</BrutalistButton>
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-2xl font-bold">Button Sizes</h2>
        <div className="flex flex-wrap items-center gap-4">
          <BrutalistButton size="sm">Small</BrutalistButton>
          <BrutalistButton>Default</BrutalistButton>
          <BrutalistButton size="lg">Large</BrutalistButton>
          <BrutalistButton size="icon">+</BrutalistButton>
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-2xl font-bold">Shadow Variations</h2>
        <div className="flex flex-wrap gap-4">
          <BrutalistButton shadow="none">No Shadow</BrutalistButton>
          <BrutalistButton shadow="sm">Small Shadow</BrutalistButton>
          <BrutalistButton>Default Shadow</BrutalistButton>
          <BrutalistButton shadow="lg">Large Shadow</BrutalistButton>
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-2xl font-bold">Transform Variations</h2>
        <div className="flex flex-wrap gap-4">
          <BrutalistButton transform="none">No Transform</BrutalistButton>
          <BrutalistButton transform="sm">Small Transform</BrutalistButton>
          <BrutalistButton>Default Transform</BrutalistButton>
          <BrutalistButton transform="lg">Large Transform</BrutalistButton>
        </div>
      </section>

      <section className="space-y-4">
        <h2 className="text-2xl font-bold">Combined Examples</h2>
        <div className="flex flex-wrap gap-4">
          <BrutalistButton 
            variant="primary" 
            size="lg"
            shadow="lg"
            transform="lg"
          >
            Primary Action
          </BrutalistButton>
          
          <BrutalistButton 
            variant="secondary" 
            size="sm"
            shadow="sm"
            transform="sm"
          >
            Secondary Action
          </BrutalistButton>
          
          <BrutalistButton 
            variant="danger" 
            shadow="lg"
            transform="none"
          >
            Danger (No Transform)
          </BrutalistButton>
          
          <BrutalistButton 
            variant="success" 
            shadow="none"
            transform="default"
          >
            Success (No Shadow)
          </BrutalistButton>
        </div>
      </section>

      <section className="space-y-4 pb-12">
        <h2 className="text-2xl font-bold">Disabled State</h2>
        <div className="flex flex-wrap gap-4">
          <BrutalistButton disabled>Default Disabled</BrutalistButton>
          <BrutalistButton variant="primary" disabled>Primary Disabled</BrutalistButton>
          <BrutalistButton variant="danger" disabled>Danger Disabled</BrutalistButton>
        </div>
      </section>
    </div>
  );
} 